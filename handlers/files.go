package handlers

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"microserver/files"

	"github.com/gorilla/mux"
)

// Files is a handler for reading and writing files
type Files struct {
	log   *log.Logger
	store files.Storage
}

// NewFiles create a new Files handler
func NewFiles(s files.Storage, l *log.Logger) *Files {
	return &Files{store: s, log: l}
}

func (f *Files) UploadRest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	filename := vars["filename"]

	f.log.Println("API: Handle POST", "id", id, "filename", filename)

	f.saveFile(id, filename, w, r.Body)
}

func (f *Files) UploadMultipart(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024) // 128 kB - in memory data
	// everything is written to a disk
	if err != nil {
		f.log.Println("Bad request", "error", err)
		http.Error(w, "Expected multipart/form-data", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		f.log.Println("Bad request", "error", err)
		http.Error(w, "Invalid id - expected integer", http.StatusBadRequest)
		return
	}

	file, mh, err := r.FormFile("file")
	if err != nil {
		f.log.Println("Bad request", "error", err)
		http.Error(w, "Expected file", http.StatusBadRequest)
		return
	}

	f.saveFile(strconv.Itoa(id), mh.Filename, w, file)
}

func (f *Files) invalidURI(uri string, w http.ResponseWriter) {
	f.log.Println("API: Error - invalid path", "path", uri)
	http.Error(w, "Invalid file path, must be in format: /[id]/[filepath]", http.StatusBadRequest)
}

// saveFile saves the contents of a request to a file
func (f *Files) saveFile(id, path string, w http.ResponseWriter, r io.ReadCloser) {
	f.log.Println("API: Save file for the product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r)
	if err != nil {
		f.log.Println("API: Error - unable to save file", "error", err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
	}
}
