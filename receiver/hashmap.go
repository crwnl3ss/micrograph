package receiver

import (
	"fmt"
	"sync"
)

// HashmapStorage ...
type HashmapStorage struct {
	sync.Mutex
	s map[string]node
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

type GrafanaDatapoint struct {
	m  float64
	ts int64
}

type GrafanaQueryTarget struct {
	Target string
}

type GrafanaQueryResult struct {
	Target     string
	Datapoints []GrafanaDatapoint
}

func (s *HashmapStorage) GetGrafanaQuery(from, to int64, targets []GrafanaQueryTarget) []GrafanaQueryResult {
	queryes := []GrafanaQueryResult{}
	s.Lock()
	defer s.Unlock()
	for _, target := range targets {
		node, ok := s.s[target.Target]
		if !ok {
			continue
		}
		subQueryResult := GrafanaQueryResult{Target: target.Target}
		for idx := range node.t {
			if node.t[idx] <= from && node.t[idx] >= to {
				subQueryResult.Datapoints = append(subQueryResult.Datapoints, GrafanaDatapoint{m: node.m[idx], ts: node.t[idx]})
			}
		}
		queryes = append(queryes, subQueryResult)
	}
	return queryes
}

func (s *HashmapStorage) insert(pp *parsedPackage) error {
	s.Lock()
	defer s.Unlock()
	existedNode, ok := s.s[pp.rawNamespace]
	if !ok {
		s.s[pp.rawNamespace] = node{t: []int64{pp.ts}, m: []float64{pp.metric}}
		return nil
	}
	if existedNode.t[len(existedNode.t)-1] < pp.ts {
		existedNode.t = append(existedNode.t, pp.ts)
		existedNode.m = append(existedNode.m, pp.metric)
		return nil
	}
	return fmt.Errorf("binary search not implemented")
}
