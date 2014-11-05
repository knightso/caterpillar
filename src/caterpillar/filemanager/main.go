package filemanager

import (
	"appengine"
	"fmt"
	"io/ioutil"
	"net/http"
)

func init() {
	http.HandleFunc("/caterpillar/filemanager/upload", UploadHandler)
	http.HandleFunc("/caterpillar/filemanager/files", FileListViewHandler)
}

func UploadHandler(rw http.ResponseWriter, req *http.Request) {
	var maxMemory int64 = 1 * 1024 * 1024
	var formKey string = "filename"
	var bucketName string // empty means default bucket

	c := appengine.NewContext(req)

	rw.Header().Set("Content-type", "text/html")
	err := req.ParseMultipartForm(maxMemory)
	if err != nil {
		if err.Error() == "permission denied" {
			rw.WriteHeader(413)
			fmt.Fprintln(rw, "アップロード可能な容量を超えています。")
		} else {
			rw.WriteHeader(500)
			fmt.Fprintln(rw, err.Error())
		}
		return
	}

	file, fileHeader, err := req.FormFile(formKey)
	if err != nil {
		rw.WriteHeader(400)
		fmt.Fprintln(rw, err.Error())
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintln(rw, err.Error())
		return
	}

	mimeType := fileHeader.Header.Get("Content-Type")
	if len(mimeType) == 0 {
		rw.WriteHeader(500)
		fmt.Fprintln(rw, "couldn't get mime-type of file.")
		return
	}

	absFilename, err := Store(c, data, fileHeader.Filename, mimeType, bucketName)
	if err != nil {
		rw.WriteHeader(500)
		fmt.Fprintln(rw, err.Error())
		return
	}

	servingURL, err := GetServingURL(c, absFilename)
	if err != nil {
		rw.WriteHeader(500)
		fmt.Fprintln(rw, err.Error())
		return
	}

	err = StoreImage(c, servingURL.String(), fileHeader.Filename, absFilename)
	if err != nil {
		rw.WriteHeader(500)
		fmt.Fprintln(rw, err.Error())
		return
	}

	rw.WriteHeader(200)
	fmt.Fprintln(rw, `<html><head><script>
		window.parent.uploadSuccess();
	</script></head></html>`)
}

func FileListViewHandler(rw http.ResponseWriter, req *http.Request) {
	var formKey string = "cursor"

	c := appengine.NewContext(req)

	cursorStr := req.FormValue(formKey)
	images, next_cursor, err := GetImages(c, cursorStr)
	if err != nil {
		return
	}

	res, err := CreateJSONResponse(c, images, next_cursor)
	if err != nil {
		return
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Write(res)
}
