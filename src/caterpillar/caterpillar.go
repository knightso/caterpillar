package caterpillar

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"caterpillar/common"
	"caterpillar/model"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/knightso/base/gae/ds"
	"github.com/knightso/base/errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var leaves map[string]*Leaf

const STATIC_DIR = "static"
const LEAF_SUFFIX = ".leaf.html"
const DWAVE_START_REPL = `{{"{{"}}`
const DWAVE_END_REPL = `{{"}}"}}`
const AREA_CLASS = "ctpl_area"
const BLOCK_CLASS = "ctpl_block"
const MNG_ROOT_URL = "/caterpillar/static/mng/"

var pageUrlRegex = regexp.MustCompile(`caterpillar://(\d+)`)

const CTPLR_TMPL_REPL = `{{if .User}}{{template "CATERPILLAR" .}}{{end}}`

type Leaf struct {
	Name      string             `json:"name"`
	Alias     string             `json:"alias"`
	Template  *template.Template `json:"-"`
	Wormholes []*Wormhole        `json:"wormholes"`
}

type WormholeType string

const (
	PROPERTY WormholeType = "PROPERTY"
	AREA     WormholeType = "AREA"
)

type Wormhole struct {
	Name   string       `json:"name"`
	Alias  string       `json:"alias"`
	Global bool         `json:"global"`
	Type   WormholeType `json:"type"`
}

func init() {

	// set memcache
	ds.DefaultCache = true

	// TODO consider duplicate Wormhole names

	// read leaves
	leaves = make(map[string]*Leaf)

	reg := regexp.MustCompile(`\{\{|\}\}|\[\[[^]]*\]\]|</[bB][oO][dD][yY]>`)

	err := filepath.Walk(STATIC_DIR, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), LEAF_SUFFIX) {
			relativePath := strings.TrimPrefix(path, STATIC_DIR)
			leafName := strings.TrimSuffix(relativePath, LEAF_SUFFIX)
			if leafName[0] == '/' {
				leafName = leafName[1:]
			}

			var alias string // under construction
			alias = leafName

			tb, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}

			src := string(tb)

			holes := []*Wormhole{}

			holeset := make(map[string]struct{})

			src = reg.ReplaceAllStringFunc(src, func(s string) string {
				if s == "{{" {
					return DWAVE_START_REPL
				} else if s == "}}" {
					return DWAVE_END_REPL
				} else if strings.HasPrefix(s, "[[") {
					wh, rpl := parseWormhole(s)

					// TODO: validate duplicate holes
					if wh != nil {
						if _, ok := holeset[wh.Name]; !ok {
							holes = append(holes, wh)
							holeset[wh.Name] = struct{}{}
						}
					}

					return rpl
				} else {
					return CTPLR_TMPL_REPL + s
				}
			})

			t, err := template.New(leafName).Parse(CTPLR_TMPL)
			if err != nil {
				panic(err)
			}
			t, err = t.Parse(src)
			if err != nil {
				panic(err)
			}

			leaves[leafName] = &Leaf{leafName, alias, t, holes}
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	// initialize martini

	m := newMartini()

	m.Get("/caterpillar/login", login)

	m.Get("**/:id"+common.VIEW_PAGE_EXT, getViewPage(false))
	m.Get("**/:id"+common.EDIT_PAGE_EXT, getViewPage(true))
	m.Get("/", defaultPage)

	m.Get("/caterpillar/api/leaves", queryLeaves)
	m.Get("/caterpillar/api/property/:name", queryProperty)

	m.Group("/caterpillar/api/pages", func(r martini.Router) {
		r.Get("", queryPages)
		r.Get("/:key", getPage)
		r.Post("", postPage)
		r.Put("", putPage)
		r.Put("/:key", putPage)
	})

	m.Group("/caterpillar/api/blocks", func(r martini.Router) {
		r.Put("/:key", saveBlocks)
	})

	m.Group("/caterpillar/api/rootPage", func(r martini.Router) {
		r.Put("", putRootPage)
	})

	http.Handle("/", m)
}

func parseWormhole(s string) (*Wormhole, string) {
	r := regexp.MustCompile(`\[\[\s*(\.?)(\w+)\s*([^]\s]*)\s*\]\]$`)

	gr := r.FindStringSubmatch(s)
	if gr == nil {
		return nil, s
	}

	dot := gr[1]
	name := gr[2]
	alias := gr[3]

	initial, _ := utf8.DecodeRune([]byte(name))

	wh := &Wormhole{
		Name:   name,
		Alias:  alias,
		Global: unicode.IsUpper(initial),
	}

	if len(dot) == 0 {
		wh.Type = AREA
		return wh, fmt.Sprintf(`{{index .Areas "%s"}}`, name)
	} else {
		wh.Type = PROPERTY

		// TODO consider more
		if name == "name" {
			wh = nil
		}

		return wh, fmt.Sprintf(`{{index .Properties "%s"}}`, name)
	}
}

type myMartini struct {
	*martini.Martini
	martini.Router
}

func newMartini() *myMartini {
	r := martini.NewRouter()
	m := martini.New()

	m.Use(func(c martini.Context, r *http.Request, l *log.Logger) {
		ac := appengine.NewContext(r)
		gaelog := log.New(logWriter{ac}, l.Prefix(), l.Flags())
		c.Map(gaelog)
	})
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(martini.Static("public"))
	m.Use(func(c martini.Context, r *http.Request) {
		ac := appengine.NewContext(r)
		c.Map(ac)
		c.Map(user.Current(ac))
	})
	m.MapTo(r, (*martini.Route)(nil))
	m.Action(r.Handle)

	return &myMartini{m, r}
}

type logWriter struct {
	ac appengine.Context
}

func login(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	// TODO: check default page and redirect to page list without it

	var postLogin string
	if dest, ok := r.URL.Query()["dest"]; ok {
		postLogin = dest[0]
	} else {
		postLogin = MNG_ROOT_URL
	}

	loginURL, _ := user.LoginURL(c, postLogin)
	http.Redirect(w, r, loginURL, http.StatusFound)
}

func (w logWriter) Write(p []byte) (n int, err error) {
	w.ac.Debugf(string(p))
	return len(p), nil
}

func getViewPage(edit bool) interface{} {
	return func(params martini.Params, u *user.User, c appengine.Context, w http.ResponseWriter, r *http.Request) {
		renderPage(params["id"], u, c, w, r, edit)
	}
}

func renderPage(id string, u *user.User, c appengine.Context, w http.ResponseWriter, r *http.Request, edit bool) {

	var pageKey *datastore.Key

	pID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		// check alias
		paKey := model.NewPageAliasKey(c, id)
		var pa model.PageAlias
		if err := ds.Get(c, paKey, &pa); err != nil {
			handleError(c, w, err, http.StatusNotFound)
			return
		}
		pID = pa.PageKey.IntID()
		pageKey = pa.PageKey
	} else {
		pageKey = model.NewPageKey(c, pID)
	}

	var p model.Page
	if err := ds.Get(c, pageKey, &p); err != nil {
		handleError(c, w, err, http.StatusNotFound)
		return
	}

	if edit {
		if u == nil {
			loginURL, _ := user.LoginURL(c, p.URLBase()+common.EDIT_PAGE_EXT)
			http.Redirect(w, r, loginURL, http.StatusFound)
			return
		} else if !u.Admin {
			// TODO: prepare error page
			http.Redirect(w, r, "/caterpillar/error/invalidUser", http.StatusFound)
			return
		}
	}

	leaf := leaves[p.Leaf]
	if leaf == nil {
		errmsg := fmt.Sprintf("leaf not found:" + p.Leaf)
		handleError(c, w, errors.New(errmsg), http.StatusNotFound)
		return
	}

	tparam := struct {
		Properties map[string]interface{}
		Areas      map[string]interface{}
		User       *user.User
		Edit       bool
		PageID     int64
		ViewURL    string
		EditURL    string
		PagesURL   string
		LogoutURL  string
	}{
		make(map[string]interface{}),
		make(map[string]interface{}),
		u,
		edit,
		pID,
		"",
		"",
		"",
		"",
	}

	if u != nil {
		tparam.PagesURL = "/caterpillar/static/mng/#/queryPage"

		purl := p.URLBase()
		tparam.ViewURL = purl + common.VIEW_PAGE_EXT
		tparam.EditURL = purl + common.EDIT_PAGE_EXT

		logoutURL, err := user.LogoutURL(c, tparam.ViewURL)
		if err != nil {
			// only log
			c.Warningf("cannot get logoutURL. err:%v", err)
		}
		tparam.LogoutURL = logoutURL
	}

	futureProps := make(map[string]<-chan func() (*model.PageProperty, error))
	futureAreas := make(map[string]<-chan func() (*model.Area, []model.Block, error))

	for _, hole := range leaf.Wormholes {
		var pkey *datastore.Key
		if hole.Global {
			pkey = model.NewGlobalPageKey(c)
		} else {
			pkey = pageKey
		}
		switch hole.Type {
		case PROPERTY:
			propkey := model.NewPagePropertyKey(c, hole.Name, pkey)
			ch := getPagePropertyAsync(c, propkey)
			futureProps[hole.Name] = ch
		case AREA:
			ch := getAreaAndBlocksAsync(c, pkey, hole.Name)
			futureAreas[hole.Name] = ch
		}
	}

	pageURLs := make(map[string]string)

	for _, hole := range leaf.Wormholes {
		switch hole.Type {
		case PROPERTY:
			prop, err := (<-futureProps[hole.Name])()
			if err == nil {
				tparam.Properties[hole.Name] = prop.Value
			} else {
				// TODO handle error
				c.Errorf("%s", err)
			}
		case AREA:
			area, blocks, err := (<-futureAreas[hole.Name])()

			if err == nil {
				areasrc := renderArea(hole.Global, area, blocks, edit)

				if !edit {
					futurePages := make(map[string]<-chan func() (*model.Page, error))

					urls := pageUrlRegex.FindAllStringSubmatch(areasrc, -1)

					for _, url := range urls {
						purl := url[0]

						if _, exists := pageURLs[purl]; exists {
							continue
						}
						if _, exists := futurePages[purl]; exists {
							continue
						}

						pageID, err := strconv.ParseInt(url[1], 10, 64)
						if err != nil {
							// TODO handle error
							c.Errorf("%s", err)
							continue
						}

						pkey := model.NewPageKey(c, pageID)
						futurePages[purl] = getPageAsync(c, pkey)
					}

					for purl, ch := range futurePages {
						page, err := (<-ch)()
						if err != nil {
							// TODO handle error
							c.Errorf("%s", err)
							continue
						}
						pageURLs[purl] = page.URLBase() + common.VIEW_PAGE_EXT
					}

					areasrc = pageUrlRegex.ReplaceAllStringFunc(areasrc, func(s string) string {
						if r, true := pageURLs[s]; true {
							return r
						}
						return s
					})
				}

				tparam.Areas[hole.Name] = template.HTML(areasrc)
			} else {
				// TODO handle error
				c.Errorf("%s", err)
			}
		}
	}

	// TODO: validate reserved page name property
	// or put some prefix?
	tparam.Properties["name"] = p.Name

	// TODO: resolve charset from somewhere
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err = leaf.Template.Execute(w, tparam); err != nil {
		handleError(c, w, err, http.StatusInternalServerError)
		return
	}

	return

}

func renderArea(global bool, area *model.Area, blocks []model.Block, edit bool) string {
	areaID := area.Key.StringID()
	// forbid underscore in the area name
	s := renderBlocks(global, areaID, blocks, edit)
	if edit {
		div := `<div id="ctpl_area:%s" class="%s">%s</div>`
		s = fmt.Sprintf(div, areaID, AREA_CLASS, s)
	}
	return s
}

func renderBlocks(global bool, areaID string, blocks []model.Block, edit bool) string {
	var s string // TODO use buffer
	for _, b := range blocks {
		s += renderBlock(global, areaID, b, edit)
	}
	return s
}

func renderBlock(global bool, areaID string, block model.Block, edit bool) string {
	bkey, s := block.Info()
	// forbid underscore in the area name
	if edit {
		div := `<div id="ctpl_block:%t:%s:%d" class="%s">%s</div>`
		s = fmt.Sprintf(div, global, areaID, bkey.IntID(), BLOCK_CLASS, s)
	}
	return s
}

func defaultPage(params martini.Params, u *user.User, c appengine.Context, w http.ResponseWriter, r *http.Request) {

	rpkey := model.NewRootPageKey(c)
	var rp model.RootPage
	if err := ds.Get(c, rpkey, &rp); err != nil {
		if errors.Root(err) == datastore.ErrNoSuchEntity {
			login(c, w, r)
		} else {
			handleError(c, w, err, http.StatusInternalServerError)
		}
		return
	}

	renderPage(strconv.FormatInt(rp.PageID, 10), u, c, w, r, false)
}

func getPagePropertyAsync(c appengine.Context, key *datastore.Key) <-chan func() (*model.PageProperty, error) {
	ch := make(chan func() (*model.PageProperty, error))
	go func() {
		var prop model.PageProperty
		if err := ds.Get(c, key, &prop); err != nil {
			ch <- func() (*model.PageProperty, error) {
				return nil, err
			}
			return
		}

		ch <- func() (*model.PageProperty, error) {
			return &prop, nil
		}
	}()
	return ch
}

func getAreaAndBlocksAsync(c appengine.Context, pkey *datastore.Key, areaName string) <-chan func() (*model.Area, []model.Block, error) {
	ch := make(chan func() (*model.Area, []model.Block, error))
	go func() {

		var area model.Area
		area.Key = model.NewAreaKey(c, areaName, pkey)

		err := datastore.RunInTransaction(c, func(c appengine.Context) error {
			err := ds.Get(c, area.Key, &area)
			if err != nil && errors.Root(err) == datastore.ErrNoSuchEntity {
				// create if not exist

				// TODO make blocks extendable & initially zero blocks
				var block model.HTMLBlock
				block.Key = model.NewHTMLBlockIncompleteKey(c, area.Key)
				block.Value = areaName + ": edit here."
				err = ds.Put(c, &block)
				if err != nil {
					return err
				}

				area.Name = areaName
				area.Blocks = []*datastore.Key{block.Key}
				err = ds.Put(c, &area)
				if err != nil {
					return err
				}
			}

			return err
		}, nil)
		if err != nil {
			ch <- func() (*model.Area, []model.Block, error) {
				return nil, nil, err
			}
			return
		}

		blockChans := []<-chan func() (*model.HTMLBlock, error){}

		for _, blockKey := range area.Blocks {
			blockCh := getHTMLBlockAsync(c, blockKey)
			blockChans = append(blockChans, blockCh)
		}

		blocks := []model.Block{}

		for _, blockCh := range blockChans {
			b, err := (<-blockCh)()
			if err != nil {
				// TODO block not found
				blocks = append(blocks, &model.HTMLBlock{Value: "block not found."})
			} else {
				blocks = append(blocks, b)
			}
		}

		ch <- func() (*model.Area, []model.Block, error) {
			return &area, blocks, nil
		}
	}()
	return ch
}

func getHTMLBlockAsync(c appengine.Context, key *datastore.Key) <-chan func() (*model.HTMLBlock, error) {
	ch := make(chan func() (*model.HTMLBlock, error))
	go func() {
		var block model.HTMLBlock
		if err := ds.Get(c, key, &block); err != nil {
			ch <- func() (*model.HTMLBlock, error) {
				return nil, err
			}
			return
		}

		ch <- func() (*model.HTMLBlock, error) {
			return &block, nil
		}
	}()
	return ch
}

func getPageAsync(c appengine.Context, key *datastore.Key) <-chan func() (*model.Page, error) {
	ch := make(chan func() (*model.Page, error))
	go func() {
		var p model.Page
		if err := ds.Get(c, key, &p); err != nil {
			ch <- func() (*model.Page, error) {
				return nil, err
			}
			return
		}

		ch <- func() (*model.Page, error) {
			return &p, nil
		}
	}()
	return ch
}

func handleError(c appengine.Context, w http.ResponseWriter, err error, code int) {
	if e, ok := err.(errors.Error); ok {
		c.Errorf(e.ErrorWithStackTrace())
	}
	http.Error(w, err.Error(), code)
}
