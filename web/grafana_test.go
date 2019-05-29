package web

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/crwnl3ss/micrograph/storage"
)

type TestGrafanaRequestRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type TestGrafanaRequestTarget struct {
	Target string `json:"target"`
}

type TestGrafanaRequest struct {
	Ranges  TestGrafanaRequestRange `json:"range"`
	Targets []QueryTarget           `json:"targets"`
}

func TestGrafanaQueryRequest(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	defer wg.Wait()
	defer cancel()
	s := storage.NewInMemoryStorage(ctx, &wg, false)
	w := httptest.NewRecorder()
	badRequest := httptest.NewRequest("POST", "/query", nil)
	handler := http.HandlerFunc(query(s))
	handler.ServeHTTP(w, badRequest)
	if w.Code != http.StatusBadRequest {
		t.Errorf("unexpected response %d", w.Code)
	}

	from := "2019-05-29T15:00:00Z"
	to := "2019-05-29T15:10:00Z"
	qr1, err := json.Marshal(&TestGrafanaRequest{
		Ranges: TestGrafanaRequestRange{
			From: from,
			To:   to,
		},
		Targets: []QueryTarget{
			QueryTarget{Target: "a"},
			QueryTarget{Target: "a.b"},
		},
	})
	if err != nil {
		t.Errorf(err.Error())
	}
	w = httptest.NewRecorder()
	emptyResponse := httptest.NewRequest("POST", "/query", bytes.NewReader(qr1))
	handler.ServeHTTP(w, emptyResponse)
	if w.Code != http.StatusOK {
		t.Errorf("unexpected response code %d", w.Code)
	}
	if w.Body.String() != "[]" {
		t.Errorf("invalid response, expect `[]` got: %s", w.Body.String())
	}

	// real tests with datapoint, query ranges:
	// 2019-05-29T15:00:00Z <-> 2019-05-29T15:10:00Z
	s.InsertDataPoint("a", &storage.DataPoint{TS: time.Date(2019, 5, 29, 14, 59, 0, 0, time.UTC).Unix(), Data: 0.01})
	s.InsertDataPoint("a", &storage.DataPoint{TS: time.Date(2019, 5, 29, 15, 0, 0, 0, time.UTC).Unix(), Data: 0.02})
	s.InsertDataPoint("a", &storage.DataPoint{TS: time.Date(2019, 5, 29, 15, 1, 0, 0, time.UTC).Unix(), Data: 0.03})
	s.InsertDataPoint("a", &storage.DataPoint{TS: time.Date(2019, 5, 29, 15, 0, 10, 0, time.UTC).Unix(), Data: 0.04})
	s.InsertDataPoint("a", &storage.DataPoint{TS: time.Date(2019, 5, 29, 15, 9, 59, 0, time.UTC).Unix(), Data: 0.05})
	s.InsertDataPoint("a", &storage.DataPoint{TS: time.Date(2019, 5, 29, 15, 9, 59, 0, time.UTC).Unix(), Data: 0.06})
	s.InsertDataPoint("a", &storage.DataPoint{TS: time.Date(2019, 5, 29, 15, 10, 1, 0, time.UTC).Unix(), Data: 0.07})
	s.InsertDataPoint("a.b", &storage.DataPoint{TS: time.Date(2019, 5, 29, 15, 55, 55, 0, time.UTC).Unix(), Data: 1})

	w = httptest.NewRecorder()
	realRequest := httptest.NewRequest("POST", "/query", bytes.NewReader(qr1))
	handler.ServeHTTP(w, realRequest)
	if w.Body.String() != `[{"Target":"a","Datapoints":[[0.02,1559142000000],[0.04,1559142010000],[0.03,1559142060000],[0.06,1559142599000],[0.05,1559142599000]]},{"Target":"a.b","Datapoints":[]}]` {
		t.Errorf("unexpected result %s", w.Body.String())
	}
}
