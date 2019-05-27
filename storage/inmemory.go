package storage

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

// HashmapStorage is a simple inmemory storage
type HashmapStorage struct {
	sync.Mutex
	s                map[string]DataPoints
	snapshotFilePath string
	snapshotEnable   bool
}

// NewInMemoryStorage ...
func NewInMemoryStorage(ctx context.Context, wg *sync.WaitGroup) Storage {
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

// Close creates snapshot of inmemory storage state
// This process could be skipped with `snapshotEnable=false` flag
func (s *HashmapStorage) Close() error {
	if !s.snapshotEnable {
		log.Println("snapshot dump disabled")
		return nil
	}
	s.Lock()
	defer s.Unlock()
	log.Printf("try to create %s", s.snapshotFilePath)
	fd, err := os.Create(s.snapshotFilePath)
	if err != nil {
		return err
	}
	bSnapshot, err := json.Marshal(s.s)
	if err != nil {
		return err
	}
	n, err := fd.Write(bSnapshot)
	if err != nil {
		return err
	}
	log.Printf("snapshot created, size %d", n)
	return nil
}

// GetGrafanaTargets ...
func (s *HashmapStorage) GetGrafanaTargets() []string {
	s.Lock()
	targets := []string{}
	for k := range s.s {
		log.Println(k)
		targets = append(targets, k)
	}
	s.Unlock()
	log.Println(targets, s.s)
	return targets
}

// GrafanaQueryResult ...
type GrafanaQueryResult struct {
	Target     string          `json:"target"`
	DataPoints [][]interface{} `json:"datapoints"`
}

// GetGrafanaQuery ...
func (s *HashmapStorage) GetGrafanaQuery(from, to int64, targets []string) []GrafanaQueryResult {
	queryes := []GrafanaQueryResult{}
	s.Lock()
	defer s.Unlock()
	for _, target := range targets {
		datapoints, ok := s.s[target]
		if !ok {
			continue
		}
		subQueryResult := GrafanaQueryResult{Target: target, DataPoints: [][]interface{}{}}
		for idx := range datapoints {
			idxDP := datapoints[idx]
			if idxDP.TS >= from && idxDP.TS <= to {
				subQueryResult.DataPoints = append(subQueryResult.DataPoints, []interface{}{idxDP.Data, idxDP.TS * 1000})
			}
		}
		queryes = append(queryes, subQueryResult)
	}
	return queryes
}

// InsertDataPoint add passed DataPoint into target's timeserease data
func (s *HashmapStorage) InsertDataPoint(target string, dp *DataPoint) error {
	s.Lock()
	defer s.Unlock()
	datapoints, ok := s.s[target]
	// there is no datapoints for passed target yet, just create new one
	if !ok {
		log.Println(s.s[target])
		log.Printf("Create new target: %s", target)
		s.s[target] = DataPoints{dp}
		return nil
	}
	// Datapoins ordered by `ts`, try to add new one at the end
	if datapoints[len(datapoints)-1].TS < dp.TS {
		log.Printf("insert %v into target %s", dp, target)
		datapoints = append(datapoints, dp)
		s.s[target] = datapoints
		return nil
	}
	// otherwise use binary search to find index
	idx := sort.Search(len(datapoints), func(i int) bool {
		return datapoints[i].TS >= dp.TS
	})
	datapoints = append(datapoints, &DataPoint{})
	copy(datapoints[idx+1:], datapoints[idx:])
	datapoints[idx] = dp
	s.s[target] = datapoints
	return nil
}
