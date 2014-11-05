package caterpillar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"unicode"
	"unicode/utf8"

	"appengine"
	"appengine/datastore"
	"appengine/delay"
	"caterpillar/model"
	"code.google.com/p/go-uuid/uuid"
	"github.com/go-martini/martini"
	"github.com/knightso/base/errors"
	"github.com/knightso/base/gae/ds"
)

var (
	ErrPageAliasAlreadyExists = errors.New("caterpillar: Page alias already exists")
)

func queryLeaves(w http.ResponseWriter) {
	// TODO: cannot marshall array?
	responseJson := "["
	index := 1
	for _, value := range leaves {
		if index != 1 {
			responseJson += ","
		}

		b, err := json.Marshal(*value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		whs := string(b)
		responseJson += whs
		index++
	}
	responseJson += "]"

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := w.Write([]byte(responseJson))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func queryProperty(params martini.Params, w http.ResponseWriter, r *http.Request) {

	pn := params["name"]
	if pn != "" {
		// get property from datastore
		c := appengine.NewContext(r)
		gpageKey := model.NewGlobalPageKey(c)
		gpropKey := model.NewPagePropertyKey(c, pn, gpageKey)

		var gpp model.PageProperty
		err := ds.Get(c, gpropKey, &gpp)
		if err != nil && errors.Root(err) != datastore.ErrNoSuchEntity {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(gpp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseJson := string(b)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, err = w.Write([]byte(responseJson))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "error: property name not found.", http.StatusInternalServerError)
		return
	}

	return
}

func getPage(params martini.Params, w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "error: illegal access.", http.StatusInternalServerError)
		return
	}

	c := appengine.NewContext(r)
	keyIDStr := params["key"]
	if keyIDStr != "" {
		// get page from datastore
		intID, err := strconv.ParseInt(keyIDStr, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		key := model.NewPageKey(c, intID)

		var p model.Page
		if err := ds.Get(c, key, &p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		props, err := getPageProperties(c, key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rpkey := model.NewRootPageKey(c)

		var rp model.RootPage
		err = ds.Get(c, rpkey, &rp)
		if err != nil && errors.Root(err) != datastore.ErrNoSuchEntity {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonPage := model.NewJsonPage(&p, props, intID == rp.PageID)

		b, err := json.Marshal(jsonPage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pageJson := string(b)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, err = w.Write([]byte(pageJson))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "error: key string not found.", http.StatusInternalServerError)
		return
	}

	return
}

func queryPages(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	rpkey := model.NewRootPageKey(c)

	var rp model.RootPage
	err := ds.Get(c, rpkey, &rp)
	if err != nil && errors.Root(err) != datastore.ErrNoSuchEntity {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pages, err := getPages(c)
	if err != nil {
		c.Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonPages := make([]*model.JsonPage, len(pages))

	for i, p := range pages {
		jsonPages[i] = model.NewJsonPage(p, nil, p.Key.IntID() == rp.PageID)
	}

	b, err := json.Marshal(jsonPages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func postPage(params martini.Params, w http.ResponseWriter, r *http.Request) {
	var jp model.JsonPage
	if err := getRequestJson(w, r, &jp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	keyIDStr := params["key"]
	if keyIDStr != "" {
		http.Error(w, fmt.Sprintf("register page did not need pageID=%s", keyIDStr), http.StatusInternalServerError)
		return
	}

	c := appengine.NewContext(r)
	pageKey, err := model.GeneratePageID(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := savePage(c, pageKey, jp.Page, jp.Properties); err != nil {
		if err == ErrPageAliasAlreadyExists {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func putPage(c appengine.Context, params martini.Params, w http.ResponseWriter, r *http.Request) {
	var jp model.JsonPage
	if err := getRequestJson(w, r, &jp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	keyIDStr := params["key"]
	if keyIDStr == "" {
		http.Error(w, "pageID not found.", http.StatusBadRequest)
		return
	}

	intID, err := strconv.ParseInt(keyIDStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pageKey := model.NewPageKey(c, intID)
	if err := savePage(c, pageKey, jp.Page, jp.Properties); err != nil {
		if err == ErrPageAliasAlreadyExists {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func putRootPage(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	var rp model.RootPage
	if err := getRequestJson(w, r, &rp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if rp.PageID == 0 {
		http.Error(w, "pageID not found.", http.StatusBadRequest)
		return
	}

	rootPageKey := model.NewRootPageKey(c)

	err := datastore.RunInTransaction(c, func(c appengine.Context) error {
		var rp2put model.RootPage
		if err := ds.Get(c, rootPageKey, &rp2put); err != nil {
			if errors.Root(err) != datastore.ErrNoSuchEntity {
				return err
			} else {
				rp2put = model.RootPage{}
				rp2put.Key = rootPageKey
			}
		}

		if rp2put.PageID != rp.PageID {
			rp2put.PageID = rp.PageID
			ds.Put(c, &rp2put)
		}

		return nil
	}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func saveBlocks(params martini.Params, w http.ResponseWriter, r *http.Request) {

	keyIDStr := params["key"]
	if keyIDStr == "" {
		http.Error(w, "the pageID not found.", http.StatusInternalServerError)
		return
	}

	var b map[string]interface{}
	if err := getRequestJson(w, r, &b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c := appengine.NewContext(r)
	intID, err := strconv.ParseInt(keyIDStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pageKey := model.NewPageKey(c, intID)

	var p model.Page
	if err := ds.Get(c, pageKey, &p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = datastore.RunInTransaction(c, func(c appengine.Context) error {
		// TODO make async
		for id, value := range b {
			c.Infof("saving block. id:%s, value:%s", id, value)
			r := regexp.MustCompile(`^ctpl_block:(true|false):(\w+):(\d+)$`)
			gr := r.FindStringSubmatch(id)
			if gr == nil {
				return errors.New("illegal block id:" + id)
			}

			global, err := strconv.ParseBool(gr[1])
			if err != nil {
				return err
			}

			areaID := gr[2]
			strBlockID := gr[3]
			blockID, err := strconv.ParseInt(strBlockID, 10, 64)
			if err != nil {
				return err
			}

			var pkey *datastore.Key
			if global {
				pkey = model.NewGlobalPageKey(c)
			} else {
				pkey = pageKey
			}
			akey := model.NewAreaKey(c, areaID, pkey)
			bkey := model.NewHTMLBlockKey(c, blockID, akey)

			var block model.HTMLBlock
			if err := ds.Get(c, bkey, &block); err != nil {
				return errors.WrapOr(err)
			}

			var ok bool
			block.Value, ok = value.(string)
			if !ok {
				return errors.New(
					fmt.Sprintf("illegal block value type :%T", value))
			}

			if err = ds.Put(c, &block); err != nil {
				return errors.WrapOr(err)
			}

			// save backup entity when HTMLBlock saved.
			blocks := []*model.HTMLBlock{&block}
			backupHTMLBlockFunc.Call(c, uuid.New(), blocks)
		}
		return nil
	}, &datastore.TransactionOptions{XG: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var backupHTMLBlockFunc = delay.Func("backup", func(c appengine.Context, uuid string, blocks []*model.HTMLBlock) error {
	for i, block := range blocks {
		backupKey := model.NewHTMLBlockBackupKey(c, fmt.Sprintf("%s-%02d", uuid, i))
		var backup model.HTMLBlockBackup
		backup.HTMLBlockKey = block.Key
		backup.HTMLBlock = *block
		// without nds cache
		if _, err := datastore.Put(c, backupKey, &backup); err != nil {
			c.Errorf(err.Error())
			return err
		}
	}
	return nil
})

func getRequestJson(w http.ResponseWriter, r *http.Request, p interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(p)
	if err != nil {
		return err
	}

	return nil
}

func savePage(c appengine.Context, pageKey *datastore.Key, p *model.Page, props map[string]string) error {

	err := datastore.RunInTransaction(c, func(c appengine.Context) error {
		var p2 model.Page
		if !pageKey.Incomplete() {
			if err := ds.Get(c, pageKey, &p2); err != nil {
				if errors.WrapOr(err) != datastore.ErrNoSuchEntity {
					return errors.WrapOr(err)
				} else {
					p2 = model.Page{}
					p2.Key = pageKey
				}
			}
		}

		oldAlias := p2.Alias

		p2.Name = p.Name
		p2.Alias = p.Alias
		p2.Leaf = p.Leaf

		aliasChanged := oldAlias != p2.Alias

		if err := ds.Put(c, &p2); err != nil {
			return errors.WrapOr(err)
		}

		if aliasChanged && p2.Alias != "" {
			newAliasKey := model.NewPageAliasKey(c, p2.Alias)

			var pa model.PageAlias
			if err := ds.Get(c, newAliasKey, &pa); err == nil {
				return ErrPageAliasAlreadyExists
			} else {
				newAlias := model.PageAlias{}
				newAlias.Key = newAliasKey
				newAlias.PageKey = pageKey

				if err := ds.Put(c, &newAlias); err != nil {
					return errors.WrapOr(err)
				}
			}
		}

		if aliasChanged && oldAlias != "" {
			oldAliasKey := model.NewPageAliasKey(c, oldAlias)

			err := ds.Delete(c, oldAliasKey)
			if err != nil {
				return err
			}
		}

		// TODO: put multi
		for name, value := range props {
			initial, _ := utf8.DecodeRune([]byte(name))
			global := unicode.IsUpper(initial)

			var pkey *datastore.Key
			if global {
				pkey = model.NewGlobalPageKey(c)
			} else {
				pkey = pageKey
			}
			propKey := model.NewPagePropertyKey(c, name, pkey)

			var prop model.PageProperty
			err := ds.Get(c, propKey, &prop)
			if err == nil {
				prop.Key = propKey
				prop.Value = value
			}

			if err != nil {
				if errors.Root(err) != datastore.ErrNoSuchEntity {
					return errors.WrapOr(err)
				} else {
					prop = model.PageProperty{}
					prop.Key = propKey
					prop.Value = value
				}
			}

			if err = ds.Put(c, &prop); err != nil {
				return errors.WrapOr(err)
			}
		}

		return nil
	}, &datastore.TransactionOptions{XG: true})

	if err != nil {
		return err
	}
	return nil
}

func getPageProperties(c appengine.Context, pageKey *datastore.Key) ([]*model.PageProperty, error) {
	q := datastore.NewQuery(model.KIND_PAGE_PROPERTY).Ancestor(pageKey)

	var props []*model.PageProperty
	if err := ds.ExecuteQuery(c, q, &props); err != nil {
		return nil, errors.WrapOr(err)
	}

	return props, nil
}

func getPages(c appengine.Context) ([]*model.Page, error) {
	q := datastore.NewQuery(model.KIND_PAGE)

	var pages []*model.Page
	if err := ds.ExecuteQuery(c, q, &pages); err != nil {
		return nil, errors.WrapOr(err)
	}

	return pages, nil
}
