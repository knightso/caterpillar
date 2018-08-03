package filemanager

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/knightso/base/errors"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
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

	client, err := storage.NewClient(c)
	if err != nil {
		log.Errorf(c, "failed to create storage client: %v", err)
		return "", errors.WrapOr(err)
	}
	defer client.Close()

	wc := client.Bucket(bucketName).Object(fileName).NewWriter(c)
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
