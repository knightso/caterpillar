package filemanager

import (
	"net/url"

	"golang.org/x/net/context"
	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine/image"
	"google.golang.org/appengine/log"
)

// GCSに保存された画像の公開URLを取得します。
func GetServingURL(c context.Context, filename string) (servingURL *url.URL, err error) {
	blobKey, err := blobstore.BlobKeyForFile(c, filename)
	if err != nil {
		log.Errorf(c, "serve.go:14: %s", err.Error())
		return nil, err
	}

	opts := &image.ServingURLOptions{Secure: true}
	servingURL, err = image.ServingURL(c, blobKey, opts)
	if err != nil {
		log.Errorf(c, "serve.go:25: %s", err.Error())
		return nil, err
	}

	return servingURL, nil
}
