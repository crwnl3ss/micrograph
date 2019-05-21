package storage

import "log"

type DataPoints []DataPoint

type DataPoint struct {
	ts int64
	m  float64
}

type Storager interface {
	insertDataPoint(string, DataPoint) bool
}

// NewStorage ...
func NewStorage(t string) *HashmapStorage {
	log.Printf("Storage type: %s", t)
	return &HashmapStorage{
		s: make(map[string]DataPoints),
	}
}
