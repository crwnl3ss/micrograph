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

// Storage is a generic interface for timesereas data storages
type Storage interface {
	InsertDataPoint(string, *DataPoint) error
	GetGrafanaTargets() []string
	GetGrafanaQuery(int64, int64, []string) []GrafanaQueryResult
	Close() error
}

// GetStorageByType ...
func GetStorageByType(ctx context.Context, t string, wg *sync.WaitGroup) Storage {
	log.Printf("storage type: %s", t)
	if t == "inmemory" {
		return NewInMemoryStorage(ctx, wg)
	}
	return nil
}
