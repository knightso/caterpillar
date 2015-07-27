package filemanager

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
)

// Store binary data to GCS
func Store(c context.Context, data []byte, fileName, mimeType, bucketName string) (absFilename string, err error) {
	if bucketName == "" {
		var err error
		if bucketName, err = file.DefaultBucketName(c); err != nil {
			log.Errorf(c, "failed to get default GCS bucket name: %v", err)
			return "", err
		}
	}

	hc := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(c, storage.ScopeFullControl),
			// Note that the App Engine urlfetch service has a limit of 10MB uploads and
			// 32MB downloads.
			// See https://cloud.google.com/appengine/docs/go/urlfetch/#Go_Quotas_and_limits
			// for more information.
			Base: &urlfetch.Transport{Context: c},
		},
	}

	ctx := cloud.NewContext(appengine.AppID(c), hc)

	wc := storage.NewWriter(ctx, bucketName, fileName)
	wc.ContentType = mimeType

	if _, err := wc.Write(data); err != nil {
		log.Errorf(c, "upload file: unable to write data to bucket %q, file %q: %v", bucketName, fileName, err)
		return "", err
	}
	if err := wc.Close(); err != nil {
		log.Errorf(c, "upload file: unable to close bucket %q, file %q: %v", bucketName, fileName, err)
		return "", err
	}

	return getAbsFilename(bucketName, fileName), nil
}

func getAbsFilename(bucketName string, fileName string) string {
	return fmt.Sprintf("/gs/%s/%s", bucketName, fileName)
}
