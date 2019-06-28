package model

import (
	"context"
	"strconv"
	"strings"

	"github.com/knightso/base/gae/ds"
	"github.com/knightso/caterpillar/src/caterpillar/common"
	"google.golang.org/appengine/datastore"
)

const KIND_PAGE = "Page"

type Page struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
	Leaf  string `json:"leaf"`
	ds.Meta
}

func (p *Page) URLBase() string {
	url := "/"
	sldex := strings.LastIndex(p.Leaf, "/")
	if sldex > 0 {
		url += p.Leaf[:sldex+1]
	}
	if p.Alias != "" {
		url += p.Alias
	} else {
		url += strconv.FormatInt(p.Key.IntID(), 10)
	}
	return url
}

func NewPageKey(c context.Context, id int64) *datastore.Key {
	return datastore.NewKey(c, KIND_PAGE, "", id, nil)
}

func GeneratePageID(c context.Context) (*datastore.Key, error) {
	intID, err := ds.GenerateID(c, KIND_PAGE)
	if err != nil {
		return nil, err
	}
	return datastore.NewKey(c, KIND_PAGE, "", intID, nil), nil
}

func NewGlobalPageKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, KIND_PAGE, "Global", 0, nil)
}

type JsonPage struct {
	*Page
	ID         int64             `json:"id"`
	ViewURL    string            `json:"viewURL"`
	EditURL    string            `json:"editURL"`
	Properties map[string]string `json:"properties"`
	Root       bool              `json:"root"`
}

func NewJsonPage(p *Page, props []*PageProperty, root bool) *JsonPage {
	base := p.URLBase()
	jp := JsonPage{
		Page:       p,
		ID:         p.Key.IntID(),
		ViewURL:    base + common.VIEW_PAGE_EXT,
		EditURL:    base + common.EDIT_PAGE_EXT,
		Properties: make(map[string]string),
		Root:       root,
	}

	for _, prop := range props {
		jp.Properties[prop.Key.StringID()] = prop.Value
	}

	return &jp
}

const KIND_PAGE_ALIAS = "PageAlias"

type PageAlias struct {
	PageKey *datastore.Key `json:"pageKey"`
	ds.Meta
}

func NewPageAliasKey(c context.Context, alias string) *datastore.Key {
	return datastore.NewKey(c, KIND_PAGE_ALIAS, alias, 0, nil)
}

const KIND_PAGE_PROPERTY = "PageProperty"

type PageProperty struct {
	Value string `json:"value"`
	ds.Meta
}

func NewPagePropertyKey(c context.Context, name string, pageKey *datastore.Key) *datastore.Key {
	return datastore.NewKey(c, KIND_PAGE_PROPERTY, name, 0, pageKey)
}

const KIND_AREA = "Area"

type Area struct {
	Name   string           `json:"name"`
	Blocks []*datastore.Key `json:"blocks"`
	ds.Meta
}

func NewAreaKey(c context.Context, name string, pageKey *datastore.Key) *datastore.Key {
	return datastore.NewKey(c, KIND_AREA, name, 0, pageKey)
}

type Block interface {
	Info() (*datastore.Key, string)
}

const KIND_HTML_BLOCK = "HTMLBlock"

type HTMLBlock struct {
	Value string `datastore:",noindex" json:"value"`
	ds.Meta
}

func NewHTMLBlockKey(c context.Context, id int64, areaKey *datastore.Key) *datastore.Key {
	return datastore.NewKey(c, KIND_HTML_BLOCK, "", id, areaKey)
}

func NewHTMLBlockIncompleteKey(c context.Context, areaKey *datastore.Key) *datastore.Key {
	return datastore.NewIncompleteKey(c, KIND_HTML_BLOCK, areaKey)
}

func (b *HTMLBlock) Info() (*datastore.Key, string) {
	return b.Key, b.Value
}

const KIND_HTML_BLOCK_BACKUP = "HTMLBlockBackup"

type HTMLBlockBackup struct {
	HTMLBlockKey *datastore.Key
	HTMLBlock
}

func NewHTMLBlockBackupKey(c context.Context, uuid string) *datastore.Key {
	return datastore.NewKey(c, KIND_HTML_BLOCK_BACKUP, uuid, 0, nil)
}

const KIND_ROOT_PAGE = "RootPage"

type RootPage struct {
	PageID int64 `json:"pageId"`
	ds.Meta
}

func NewRootPageKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, KIND_ROOT_PAGE, "Root", 0, nil)
}
