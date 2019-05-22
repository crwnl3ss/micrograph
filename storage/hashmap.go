package storage

import (
	"fmt"
	"sync"
)

// HashmapStorage ...
type HashmapStorage struct {
	sync.Mutex
	s map[string]DataPoints
}

// GetGrafanaTargets ...
func (s *HashmapStorage) GetGrafanaTargets() []string {
	s.Lock()
	targets := []string{}
	for k := range s.s {
		targets = append(targets, k)
	}
	s.Unlock()
	return targets
}

// GrafanaQueryTarget ...
type GrafanaQueryTarget struct {
	Target string
}

// GrafanaQueryResult ...
type GrafanaQueryResult struct {
	Target string
	DataPoints
}

// GetGrafanaQuery ...
func (s *HashmapStorage) GetGrafanaQuery(from, to int64, targets []GrafanaQueryTarget) []GrafanaQueryResult {
	queryes := []GrafanaQueryResult{}
	s.Lock()
	defer s.Unlock()
	for _, target := range targets {
		datapoints, ok := s.s[target.Target]
		if !ok {
			continue
		}
		subQueryResult := GrafanaQueryResult{Target: target.Target}
		for idx := range datapoints {
			idxDP := datapoints[idx]
			if idxDP.TS <= from && idxDP.TS >= to {
				subQueryResult.DataPoints = append(subQueryResult.DataPoints, idxDP)
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
		s.s[target] = DataPoints{dp}
		return nil
	}
	// Datapoins ordered by `ts`, try to add new one at the end
	if datapoints[len(datapoints)-1].TS < dp.TS {
		datapoints = append(datapoints, dp)
		return nil
	}
	// otherwise use binary search for new Datapoint
	return fmt.Errorf("binary search not implemented")
}
