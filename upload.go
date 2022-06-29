package dragonSpider

import (
	"errors"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/kaliadmen/dragon_spider/v2/filesystems"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

func (ds *DragonSpider) UploadFile(w http.ResponseWriter, r *http.Request, destination, field string, fs filesystems.Fs) error {
	fileName, err := ds.getFileToUpload(w, r, field)
	if err != nil {
		return err
	}

	//uploading to a remote filesystem
	if fs != nil {
		err = fs.Put(fileName, destination)
		if err != nil {
			ds.ErrorLog.Println(err)
			return err
		}
	} else { //uploading to local filesystem
		err = os.Rename(fileName, fmt.Sprintf("%s/%s", destination, path.Base(fileName)))
		if err != nil {
			ds.ErrorLog.Println(err)
			return err
		}
	}

	defer func() {
		err := os.Remove(fileName)
		if err != nil {
			ds.ErrorLog.Println(fileName, "could not be removed")
		}
	}()

	return nil
}

func (ds *DragonSpider) getFileToUpload(w http.ResponseWriter, r *http.Request, fieldName string) (string, error) {
	if r.ContentLength > ds.config.uploads.maxUploadSize {
		return "", errors.New("upload size exceeded")
	}

	r.Body = http.MaxBytesReader(w, r.Body, ds.config.uploads.maxUploadSize)

	if err := r.ParseMultipartForm(ds.config.uploads.maxUploadSize); err != nil {
		return "", err
	}

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", err
	}

	if header.Size > ds.config.uploads.maxUploadSize {
		return "", errors.New("upload size exceeded")
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			_ = errors.New("file could not be closed")
			return
		}
	}(file)

	//reads the first ~512 bytes of the file
	mimeType, err := mimetype.DetectReader(file)
	if err != nil {
		return "", err
	}
	//go back to start of file
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	if !inSlice(ds.config.uploads.allowedMimeTypes, mimeType.String()) {
		return "", errors.New("invalid file type uploaded")
	}

	dest, err := os.Create(fmt.Sprintf("./tmp/%s", header.Filename))
	if err != nil {
		return "", err
	}

	defer func(dest *os.File) {
		err := dest.Close()
		if err != nil {
			_ = errors.New("destination file could not be closed")
			return
		}
	}(dest)

	_, err = io.Copy(dest, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("./tmp/%s", header.Filename), nil
}

func inSlice(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}

	return false
}
