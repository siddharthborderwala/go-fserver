package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	maxFileSize int // maximum number of bytes for files
	basePath    string
}

// NewLocal creates a ew local filesystem with the given base path
// basepath is the base directory to save files to
// maxSize is the max number of bytes that a file can be
func NewLocal(basePath string, maxSize int) (*Local, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	return &Local{basePath: p, maxFileSize: maxSize}, nil
}

// Save the contens of the Writer to the given path
// oath is a relative path, basePath will be appended
func (l *Local) Save(path string, contents io.Reader) error {
	// get the full path for the file
	fp := l.fullPath(path)

	// get the directory and make sure it exists
	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to create directory: %w", err)
	}

	// if the file exists delete it
	_, err = os.Stat(fp)
	if err == nil {
		err = os.Remove(fp)
		if err != nil {
			return fmt.Errorf("unable to delete existing file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		// if this is anything other than a not exists error
		return fmt.Errorf("unable to get file info: %w", err)
	}

	// create a new file at the path
	f, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer f.Close()

	// write the contens to the new file
	// ensure that we are not writing greater than maxBytes
	_, err = io.Copy(f, contents)
	if err != nil {
		return fmt.Errorf("unable to write to file: %w", err)
	}

	return nil
}

// Get the files ar the given path and return a Reader
// the calling function is responsible for closing the reader
func (l *Local) Get(path string) (*os.File, error) {
	// get the full path for the file
	fp := l.fullPath(path)

	// open the file
	f, err := os.Open(fp)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	return f, nil
}

// returns the absolute path
func (l *Local) fullPath(path string) string {
	return filepath.Join(l.basePath, path)
}
