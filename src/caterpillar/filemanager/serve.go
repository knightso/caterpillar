package filemanager

import (
	"net/url"

	"golang.org/x/net/context"
	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine/image"
	"google.golang.org/appengine/log"
)

// GetServingURL gets URL to serve GCS file to users.
func GetServingURL(c context.Context, filename string) (servingURL *url.URL, err error) {

	blobKey, err := blobstore.BlobKeyForFile(c, filename)
	if err != nil {
		log.Errorf(c, "error for file %s:  %s", filename, err.Error())
		return nil, err
	}

	opts := &image.ServingURLOptions{Secure: true}
	servingURL, err = image.ServingURL(c, blobKey, opts)
	if err != nil {
		log.Errorf(c, "failed to get serving URL: %s", err.Error())
		return nil, err
	}

	return servingURL, nil
}
