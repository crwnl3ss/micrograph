package storage

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestInsertAndRetrieveDataPoint(t *testing.T) {
	wg := new(sync.WaitGroup)
	ctx, ctxCancelFn := context.WithCancel(context.Background())
	defer wg.Wait()
	defer ctxCancelFn()
	s := NewInMemoryStorage(ctx, wg, false)
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
	secondT := "d.c.b.a"
	secondTfirstDP := &DataPoint{
		TS:   time.Date(2019, 1, 1, 10, 0, 15, 0, time.UTC).Unix(),
		Data: 0.01,
	}
	if err := s.InsertDataPoint(secondT, secondTfirstDP); err != nil {
		t.Errorf("could not insert %v", secondTfirstDP)
	}
	q := s.RangeQuery(1, 9999999999, []string{firstT})
	if len(q) != 1 {
		t.Errorf("`%s` target: wromg number of datapoints: %d", firstT, len(q))
	}
	expect := RangeQueryResult{firstT, []*DataPoint{firstTfirstDP, firstTsecondDP}}
	if q[0].Target != expect.Target {
		t.Errorf("invalid target, exppect: %s, got: %s", expect.Target, q[0].Target)
	}
	if len(q[0].DataPoints) != len(expect.DataPoints) {
		t.Errorf("invalid count of datapoints expect: %d, got: %d", len(expect.DataPoints), len(q[0].DataPoints))
	}

	if expect.DataPoints[0] != q[0].DataPoints[0] || expect.DataPoints[1] != q[0].DataPoints[1] {
		t.Errorf("invalid datapoint for target %s", firstT)
	}
	// test against second target
	q = s.RangeQuery(1, 9999999999, []string{secondT})
	if len(q) != 1 {
		t.Errorf("wrong number of targets exppect: %d got %d", 1, len(q))
	}
	if q[0].Target != secondT {
		t.Errorf("invalid target expect: %s got: %s", secondT, q[0].Target)
	}

	if len(q[0].DataPoints) != 1 || q[0].DataPoints[0] != secondTfirstDP {
		t.Errorf("invalid datapoints query result: expect %+v got %+v", secondTfirstDP, q[0].DataPoints)
	}
	// add some DataPoints to first targer again
	firstTthirdDP := &DataPoint{
		TS:   time.Date(2019, 1, 1, 10, 0, 40, 0, time.UTC).Unix(),
		Data: 0.01,
	}
	if err := s.InsertDataPoint(firstT, firstTthirdDP); err != nil {
		t.Errorf(err.Error())
	}
	// test insert using binary search
	firstTsecondAndHalfDP := &DataPoint{
		TS:   time.Date(2019, 1, 1, 10, 0, 20, 0, time.UTC).Unix(),
		Data: 2.5,
	}
	if err := s.InsertDataPoint(firstT, firstTsecondAndHalfDP); err != nil {
		t.Fatalf(err.Error())
	}
	q = s.RangeQuery(1, 9999999999, []string{firstT})
	if len(q[0].DataPoints) != 4 || q[0].DataPoints[0] != firstTfirstDP || q[0].DataPoints[1] != firstTsecondDP || q[0].DataPoints[2] != firstTsecondAndHalfDP || q[0].DataPoints[3] != firstTthirdDP {
		t.Errorf("target: %s invalid datapoins count or order", firstT)
	}

	firstTsecondAndThreeQUotersDP := &DataPoint{
		TS:   time.Date(2019, 1, 1, 10, 0, 37, 0, time.UTC).Unix(),
		Data: 2.5,
	}
	if err := s.InsertDataPoint(firstT, firstTsecondAndThreeQUotersDP); err != nil {
		t.Fatalf(err.Error())
	}
	q = s.RangeQuery(1, 9999999999, []string{firstT})
	if len(q[0].DataPoints) != 5 || q[0].DataPoints[0] != firstTfirstDP || q[0].DataPoints[1] != firstTsecondDP || q[0].DataPoints[2] != firstTsecondAndHalfDP || q[0].DataPoints[3] != firstTsecondAndThreeQUotersDP || q[0].DataPoints[4] != firstTthirdDP {
		t.Errorf("target: %s invalid datapoins count or order %+v", firstT, q[0].DataPoints)
	}

}

func TestGetKeys(t *testing.T) {
	wg := new(sync.WaitGroup)
	ctx, ctxCancelFn := context.WithCancel(context.Background())
	defer wg.Wait()
	defer ctxCancelFn()

	s := NewInMemoryStorage(ctx, wg, false)
	s.InsertDataPoint("a", &DataPoint{2, 2.0})
	s.InsertDataPoint("a", &DataPoint{4, 2.2})
	s.InsertDataPoint("a.b", &DataPoint{1, 77.7})
	s.InsertDataPoint("a.b", &DataPoint{2, 2.2})
	s.InsertDataPoint("a.b", &DataPoint{9, 9.02})
	get := s.GetKeys()
	if get[0] == "a" && get[1] != "a.b" || get[0] == "a.b" && get[1] != "a" {
		t.Errorf("Unexpected `GetGrafanaTargets` result: %s", get)
	}
}
