package filemanager

import (
	"appengine"
	"appengine/blobstore"
	"appengine/image"
	"net/url"
)

// GCSに保存された画像の公開URLを取得します。
func GetServingURL(c appengine.Context, filename string) (servingURL *url.URL, err error) {
	blobKey, err := blobstore.BlobKeyForFile(c, filename)
	if err != nil {
		c.Errorf("serve.go:14: %s", err.Error())
		return nil, err
	}

	opts := &image.ServingURLOptions{Secure: true}
	servingURL, err = image.ServingURL(c, blobKey, opts)
	if err != nil {
		c.Errorf("serve.go:25: %s", err.Error())
		return nil, err
	}

	return servingURL, nil
}
