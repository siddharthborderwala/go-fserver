package files

import "io"

// Storage defines the behavior for the file operations
// Implementations may be of a time local disk, or cloud storage, etc
type Storage interface {
	Save(path string, file io.Reader) error
}
