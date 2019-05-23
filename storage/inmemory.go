package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

// HashmapStorage ...
type HashmapStorage struct {
	sync.Mutex
	s                map[string]DataPoints
	snapshotFilePath string
}

// Close creates snapshot of curret inmemory storage (at least tryes...)
func (s *HashmapStorage) Close() error {
	s.Lock()
	defer s.Unlock()
	log.Println("Creating snapshot...")
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
				subQueryResult.DataPoints = append(subQueryResult.DataPoints, []interface{}{idxDP.Data, idxDP.TS})
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
		return nil
	}
	// otherwise use binary search for new Datapoint
	return fmt.Errorf("binary search not implemented")
}
