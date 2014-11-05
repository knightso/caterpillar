package filemanager

import (
	"appengine"
	"appengine/file"
)

// アップロードされた画像をGCSに保存します。
func Store(c appengine.Context, data []byte, filename, mimeType, bucketName string) (absFilename string, err error) {

	opts := &file.CreateOptions{
		MIMEType:   mimeType,
		BucketName: bucketName,
	}
	wc, absFilename, err := file.Create(c, filename, opts)
	if err != nil {
		c.Errorf("gcs.go:23")
		return "", err
	}
	defer wc.Close()

	_, err = wc.Write(data)
	if err != nil {
		c.Errorf("gcs.go:30")
		return "", err
	}

	return absFilename, nil
}
