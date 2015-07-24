package filemanager

import (
	"golang.org/x/net/context"
	//"google.golang.org/appengine/file"
	//"google.golang.org/appengine/log"
)

// アップロードされた画像をGCSに保存します。
func Store(c context.Context, data []byte, filename, mimeType, bucketName string) (absFilename string, err error) {

	// TODO: need fix File API.
	/*opts := &file.CreateOptions{
		MIMEType:   mimeType,
		BucketName: bucketName,
	}
	wc, absFilename, err := file.Create(c, filename, opts)
	if err != nil {
		log.Errorf("gcs.go:23")
		return "", err
	}
	defer wc.Close()

	_, err = wc.Write(data)
	if err != nil {
		log.Errorf("gcs.go:30")
		return "", err
	}

	return absFilename, nil*/
	return "", nil
}
