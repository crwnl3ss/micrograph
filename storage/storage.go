package storage

import (
	"context"
	"log"
	"sync"
)

// DataPoints stores ordered by timestamp collection of Datapoint
type DataPoints []*DataPoint

// DataPoint stores timeserias data: timestamp and metric
type DataPoint struct {
	TS   int64
	Data float64
}

// RangeQueryResult ...
// TODO: serialize for grafana correctrly
type RangeQueryResult struct {
	Target     string       `json:"target"`
	DataPoints []*DataPoint `json:"datapoints"`
}

// Storage is a generic interface for timesereas data storages
type Storage interface {
	InsertDataPoint(string, *DataPoint) error
	GetKeys() []string
	RangeQuery(int64, int64, []string) []RangeQueryResult
	Close() error
}

// GetStorageByType ...
func GetStorageByType(ctx context.Context, t string, wg *sync.WaitGroup) Storage {
	log.Printf("storage type: %s", t)
	if t == "inmemory" {
		return NewInMemoryStorage(ctx, wg, true)
	}
	return nil
}
