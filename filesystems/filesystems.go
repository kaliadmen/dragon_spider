package filesystems

import "time"

//Fs is an interface for filesystems.
//All functions in this interface must exist in order to be implemented.
type Fs interface {
	Put(fileName, directory string) error
	Get(destination string, items ...string) error
	List(prefix string) ([]Listing, error)
	Delete(itemsToDelete []string) (bool, error)
}

//Listing represents a file on a remote file system
type Listing struct {
	Etag         string
	LastModified time.Time
	Key          string
	Size         float64
	IsDir        bool
}

func ToMegabytes(bytes float64) float64 {
	return (bytes / 1024) / 1024
}
