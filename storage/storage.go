package storage

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// DataPoints stores ordered by timestamp collection of Datapoint
type DataPoints []*DataPoint

// DataPoint stores timeserias data: timestamp and metric
type DataPoint struct {
	TS   int64
	Data float64
}

// Storager is a generic interface for timesereas data storages
type Storager interface {
	InsertDataPoint(string, DataPoint) error
	Close() error
}

// NewStorage creates new inmemory storage
func NewStorage(ctx context.Context, t string, wg *sync.WaitGroup) *HashmapStorage {
	log.Printf("Storage type: %s", t)
	wg.Add(1)
	s := &HashmapStorage{
		s:                make(map[string]DataPoints),
		snapshotFilePath: "./mg.snapshot",
	}
	s.Lock()
	defer s.Unlock()
	go func() {
		<-ctx.Done()
		if err := s.Close(); err != nil {
			log.Println(err)
		}
		wg.Done()
	}()
	log.Println("Looking for snapshots...")
	fd, err := os.Open(s.snapshotFilePath)
	if err != nil {
		log.Printf("Could not load snapshot. Reason: %s", err.Error())
		log.Println("inmemory storage will be empty")
		return s
	}
	bSnapshot, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Printf("Could not load snapshot. Reason: %s", err.Error())
		log.Println("inmemory storage empty")
		return s
	}
	if err := json.Unmarshal(bSnapshot, &s.s); err != nil {
		log.Printf("could not deserialize %s file. Reason %s", s.snapshotFilePath, err)
		return s
	}
	log.Println("Snapshot succsessful load")
	return s
}
