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

func (s *HashmapStorage) GetGrafanaTarggets() []string {
	s.Lock()
	targets := []string{}
	for k := range s.s {
		targets = append(targets, k)
	}
	s.Unlock()
	return targets
}

type GrafanaQueryTarget struct {
	Target string
}

type GrafanaQueryResult struct {
	Target string
	DataPoints
}

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
			if idxDP.ts <= from && idxDP.ts >= to {
				subQueryResult.DataPoints = append(subQueryResult.DataPoints, DataPoint{})
			}
		}
		queryes = append(queryes, subQueryResult)
	}
	return queryes
}

// insert passed Datapoint into slice of Datapoints of passed target
func (s *HashmapStorage) insertDataPoint(target string, dp DataPoint) error {
	s.Lock()
	defer s.Unlock()
	datapoints, ok := s.s[target]
	// there is no datapoints for passed target yet, just create new one
	if !ok {
		s.s[target] = DataPoints{dp}
		return nil
	}
	// Datapoins ordered by `ts`, try to add new one at the end
	if datapoints[len(datapoints)-1].ts < dp.ts {
		datapoints = append(datapoints, dp)
		return nil
	}
	// otherwise use binary search for new Datapoint
	return fmt.Errorf("binary search not implemented")
}
