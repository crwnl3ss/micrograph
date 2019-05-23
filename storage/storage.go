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
	log.Printf("storage type: %s", t)
	wg.Add(1)
	s := &HashmapStorage{
		s:                make(map[string]DataPoints),
		snapshotFilePath: "./mg.snapshot",
		snapshotEnable:   true,
	}
	s.Lock()
	defer s.Unlock()
	go func() {
		<-ctx.Done()
		if err := s.Close(); err != nil {
			log.Printf("could not properly close storage, reason: %s", err)
		}
		wg.Done()
	}()
	if !s.snapshotEnable {
		log.Println("snapshot load/dump disabled")
		return s
	}
	log.Println("looking for prevous snapshot...")
	fd, err := os.Open(s.snapshotFilePath)
	if err != nil {
		log.Printf("could not open %s, reason: %s", s.snapshotFilePath, err.Error())
		return s
	}
	bSnapshot, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Printf("could not read %s, reason: %s", s.snapshotFilePath, err.Error())
		return s
	}
	if err := json.Unmarshal(bSnapshot, &s.s); err != nil {
		log.Printf("could not deserialize %s file, reason %s", s.snapshotFilePath, err)
		return s
	}
	log.Println("snapshot succsessful load <3")
	return s
}
