package storage

import (
	"testing"
)

func TestGetGrafanaTargets(t *testing.T) {
	s := NewStorage("hashmap")
	s.s["a"] = DataPoints{
		&DataPoint{2, 2.0},
		&DataPoint{4, 2.2},
	}
	s.s["a.b"] = DataPoints{
		&DataPoint{1, 77.7},
		&DataPoint{2, 2.2},
		&DataPoint{9, 9.02},
	}
	get := s.GetGrafanaTargets()
	if get != []string{"a", "a.b"} || get != []string{"a.b", "a"} {
		t.Errorf("Unexpected `GetGrafanaTargets` result: %s", get)
	}
}
