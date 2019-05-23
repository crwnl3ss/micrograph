package storage

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestInsertDataPoint(t *testing.T) {
	wg := new(sync.WaitGroup)
	ctx, ctxCancelFn := context.WithCancel(context.Background())
	defer wg.Wait()
	defer ctxCancelFn()
	s := NewStorage(ctx, "hashmap", wg)
	s.snapshotEnable = false
	firstTfirstDP := &DataPoint{
		TS:   time.Date(2019, 1, 1, 10, 0, 0, 0, time.UTC).Unix(),
		Data: 0.01,
	}
	firstT := "a.a"
	firstTsecondDP := &DataPoint{
		TS:   time.Date(2019, 1, 1, 10, 0, 1, 0, time.UTC).Unix(),
		Data: 0.01,
	}
	if err := s.InsertDataPoint("a.a", firstTfirstDP); err != nil {
		t.Errorf("could not insert %v", firstTfirstDP)
	}
	if err := s.InsertDataPoint(firstT, firstTsecondDP); err != nil {
		t.Errorf("could not insert %v", firstTsecondDP)
	}
	secondTfirstDP := &DataPoint{
		TS:   time.Date(2019, 1, 1, 10, 0, 15, 0, time.UTC).Unix(),
		Data: 0.01,
	}
	if err := s.InsertDataPoint("d.c.b.a", secondTfirstDP); err != nil {
		t.Errorf("could not insert %v", secondTfirstDP)
	}
	if len(s.s[firstT]) != 2 {
		t.Errorf("`%s` target: wromg number of datapoints: %d", firstT, len(s.s["a.a"]))
	}
	if s.s[firstT][0] != firstTfirstDP || s.s[firstT][1] != firstTsecondDP {
		t.Errorf("`%s` target: invalid datapoints: %v %v", firstT, firstTfirstDP, firstTsecondDP)
	}
}

func TestGetGrafanaQuery(t *testing.T) {
	wg := new(sync.WaitGroup)
	ctx, ctxCancelFn := context.WithCancel(context.Background())
	defer wg.Wait()
	defer ctxCancelFn()
	s := NewStorage(ctx, "hashmap", wg)
	s.snapshotEnable = false
	s.s["a.b.c"] = DataPoints{
		&DataPoint{
			TS:   time.Date(2019, 1, 1, 10, 0, 0, 0, time.UTC).Unix(),
			Data: 0.2,
		},
		&DataPoint{
			TS:   time.Date(2019, 1, 1, 12, 0, 0, 0, time.UTC).Unix(),
			Data: 0.3,
		},
		&DataPoint{
			TS:   time.Date(2019, 1, 1, 13, 0, 0, 0, time.UTC).Unix(),
			Data: 0.1,
		},
	}
	s.s["target#1"] = DataPoints{
		&DataPoint{
			TS:   time.Date(2019, 1, 1, 10, 30, 0, 0, time.UTC).Unix(),
			Data: 0.2,
		},
		&DataPoint{
			TS:   time.Date(2019, 1, 1, 10, 40, 0, 0, time.UTC).Unix(),
			Data: 0.3,
		},
		&DataPoint{
			TS:   time.Date(2019, 1, 1, 10, 40, 30, 0, time.UTC).Unix(),
			Data: 0.1,
		},
	}
	res := s.GetGrafanaQuery(
		time.Date(2019, 1, 1, 9, 0, 0, 0, time.UTC).Unix(),
		time.Date(2019, 1, 1, 10, 45, 0, 0, time.UTC).Unix(),
		[]string{"target#1"},
	)
	if len(res[0].DataPoints) != 3 {
		t.Errorf("Unexpected `GetGrafanaQuery` result for target `target#1`")
	}
}

func TestGetGrafanaTargets(t *testing.T) {
	wg := new(sync.WaitGroup)
	ctx, ctxCancelFn := context.WithCancel(context.Background())
	defer wg.Wait()
	defer ctxCancelFn()

	s := NewStorage(ctx, "hashmap", wg)
	s.snapshotEnable = false
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
	if get[0] == "a" && get[1] != "a.b" || get[0] == "a.b" && get[1] != "a" {
		t.Errorf("Unexpected `GetGrafanaTargets` result: %s", get)
	}
}
