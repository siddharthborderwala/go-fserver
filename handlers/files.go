package handlers

import (
	"log"
	"net/http"
	"path/filepath"

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

// ServeHTTP is needed for http.Handler interface's implementation
func (f *Files) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	filename := vars["filename"]

	f.log.Println("API: Handle POST", "id", id, "filename", filename)

	f.saveFile(id, filename, w, r)
}

func (f *Files) invalidURI(uri string, w http.ResponseWriter) {
	f.log.Println("API: Error - invalid path", "path", uri)
	http.Error(w, "Invalid file path, must be in format: /[id]/[filepath]", http.StatusBadRequest)
}

// saveFile saves the contents of a request to a file
func (f *Files) saveFile(id, path string, w http.ResponseWriter, r *http.Request) {
	f.log.Println("API: Save file for the product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r.Body)
	if err != nil {
		f.log.Println("API: Error - unable to save file", "error", err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
	}
}
