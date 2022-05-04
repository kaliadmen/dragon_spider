package dragonSpider

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
)

func (ds *DragonSpider) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DragonSpider) WriteXML(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DragonSpider) FileDownload(w http.ResponseWriter, r *http.Request, pathToFile, filename string) error {
	fp := path.Join(pathToFile, filename)
	fileToServer := filepath.Clean(fp)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	http.ServeFile(w, r, fileToServer)
	return nil
}

func (ds *DragonSpider) ErrorStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (ds *DragonSpider) Error401(w http.ResponseWriter, r *http.Request) {
	ds.ErrorStatus(w, http.StatusUnauthorized)
}

func (ds *DragonSpider) Error403(w http.ResponseWriter, r *http.Request) {
	ds.ErrorStatus(w, http.StatusForbidden)
}

func (ds *DragonSpider) Error404(w http.ResponseWriter, r *http.Request) {
	ds.ErrorStatus(w, http.StatusNotFound)
}

func (ds *DragonSpider) Error500(w http.ResponseWriter, r *http.Request) {
	ds.ErrorStatus(w, http.StatusInternalServerError)
}
