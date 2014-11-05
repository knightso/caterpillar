package filemanager

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/knightso/base/gae/ds"
)

const (
	thumbnailsLongestSide int = 100
)

const (
	FILE  Type = "file"
	IMAGE      = "image"
)

const KIND_FILE = "File"

type Type string

type File struct {
	ServingURL   string `datastore:",noindex" json:"url"`
	FileName     string `datastore:",noindex" json:"filename"`
	GCSPath      string `datastore:",noindex" json:"-"` // gcs.go:Store()の戻り値absFilename
	Type         Type
	ThumbnailURL string `datastore:"-" json:"thumbnail"`
	ds.Meta
}

// File構造体のThumbnailURLを設定します。
// size = 0で元の大きさで表示されるようだ。1600以上でも同様のようです。
func (i *File) setThumbnailURL(size int, isCrop bool) error {
	u := i.ServingURL
	u += fmt.Sprintf("=s%d", size)
	if isCrop {
		u += "-c"
	}
	pu, err := url.Parse(u)
	if err != nil {
		return err
	}
	i.ThumbnailURL = pu.String()
	return nil
}

func GetParentKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, KIND_FILE, "/", 0, nil)
}

func NewFileKey(c appengine.Context, parentKey *datastore.Key) *datastore.Key {
	return datastore.NewIncompleteKey(c, KIND_FILE, parentKey)
}

// アップロードされた画像のメタデータをDSに保存します。
func StoreImage(c appengine.Context, servingURL, fileName, gcsPath string) error {
	e := File{
		ServingURL: servingURL,
		FileName:   fileName,
		GCSPath:    gcsPath,
		Type:       IMAGE,
	}

	e.Key = NewFileKey(c, GetParentKey(c))
	if err := ds.Put(c, &e);err != nil {
		return err
	}

	return nil
}

// 画像のメタデータ一覧をDSから取得します。
// TODO: 表示する画像数を絞る必要がないなら、Cursor必要ないかも。
func GetImages(c appengine.Context, cursorStr string) ([]File, string, error) {
	parentKey := GetParentKey(c)
	q := datastore.NewQuery(KIND_FILE).Ancestor(parentKey).Filter("Type =", IMAGE).Order("-CreatedAt")

	if len(cursorStr) != 0 {
		cursor, err := datastore.DecodeCursor(cursorStr)
		if err != nil {
			return []File{}, "", err
		}

		q = q.Start(cursor)
	}

	images := []File{}
	iter := q.Run(c)
	isNext := true
	for {
		var img File
		_, err := iter.Next(&img)
		if err == datastore.Done {
			isNext = false
			break
		}
		if err != nil {
			c.Errorf("fetching next File: %s", err.Error())
			break
		}

		err = img.setThumbnailURL(thumbnailsLongestSide, false)
		if err != nil {
			c.Errorf("%s", err.Error())
			break
		}
		images = append(images, img)
	}

	if isNext {
		next_cursor, err := iter.Cursor()
		if err != nil {
			c.Errorf("%s", err.Error())
			return []File{}, "", err
		}
		return images, next_cursor.String(), nil
	} else {
		return images, "", nil
	}
}

type ResponseData struct {
	Files  []File `json:"files"`
	Cursor string `json:"cursor"`
}

// ファイルの一覧をクライアントに返すため、FileのスライスをJSONシリアライズします。
func CreateJSONResponse(c appengine.Context, images []File, cursor string) ([]byte, error) {
	fl := &ResponseData{Files: images, Cursor: cursor}
	res, err := json.Marshal(fl)
	if err != nil {
		c.Errorf("%s", err.Error())
		return []byte{}, err
	}

	return res, nil
}
