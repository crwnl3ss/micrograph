package storage

import "log"

// DataPoints stores ordered by timestamp collection of Datapoint
type DataPoints []*DataPoint

// DataPoint stores timeserias data: timestamp and metric
type DataPoint struct {
	TS   int64
	Data float64
}

// Storager is a generic interface for timesereas data storages
type Storager interface {
	InsertDataPoint(string, DataPoint) bool
}

// NewStorage ...
func NewStorage(t string) *HashmapStorage {
	log.Printf("Storage type: %s", t)
	return &HashmapStorage{
		s: make(map[string]DataPoints),
	}
}
