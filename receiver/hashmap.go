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
	for k, _ := range s.s {
		targets = append(targets, k)
	}
	s.Unlock()
	return targets
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
