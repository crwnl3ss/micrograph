package receiver

import "log"

type node struct {
	t []int64
	m []float64
}

type storager interface {
	insert(pp *parsedPackage) bool
}

// NewStorage ...
func NewStorage(t string) *HashmapStorage {
	log.Printf("Storage type: %s", t)
	return &HashmapStorage{
		s: make(map[string]node),
	}
}
