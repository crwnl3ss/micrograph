package web

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/crwnl3ss/micrograph/storage"
)

var startTime = time.Now()

// grafana datasorce index handler
func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(string("minigraph grafana exporter v0.1")))
}

// SearchRequest ...
type SearchRequest struct {
	Target string `json:"target,omitempty"`
}

// SearchResponse ...
type SearchResponse []string

// search returns list of grafana targets
func search(s *storage.HashmapStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		sReq := &SearchRequest{}
		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(b, sReq); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println(sReq.Target)
		var sRes SearchResponse = s.GetGrafanaTargets()
		br, err := json.Marshal(sRes)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadGateway)
		}
		w.Write(br)
		w.Header().Add("Content-Type", "application/json")
	}
}

// UnmarshalJSON ...
func (qr *QueryRanges) UnmarshalJSON(b []byte) error {
	tmp := make(map[string]interface{})
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	from, err := time.Parse("2006-01-02T15:04:05Z", tmp["from"].(string))
	if err != nil {
		return err
	}
	qr.From = from.Unix()
	to, err := time.Parse("2006-01-02T15:04:05Z", tmp["to"].(string))
	qr.To = to.Unix()
	return nil
}

// QueryRanges ...
type QueryRanges struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`
}

// QueryTarget ...
type QueryTarget struct {
	Target string `json:"target"`
}

// QueryRequest ...
type QueryRequest struct {
	Range   QueryRanges   `json:"range"`
	Targets []QueryTarget `json:"targets"`
}

// query request for grafana simple json datasource
func query(s *storage.HashmapStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// in debug purpuses
			targets := s.GetGrafanaTargets()
			queries := s.GetGrafanaQuery(startTime.Unix(), time.Now().Unix(), targets)
			b, err := json.Marshal(queries)
			if err != nil {
				log.Fatalln(err)
			}
			w.Write(b)
			w.Header().Add("Content-Type", "application/json")
			return
		}
		if r.Method == "POST" {
			bReq, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatalln(err)
			}
			defer r.Body.Close()
			qr := &QueryRequest{}
			if err := json.Unmarshal(bReq, qr); err != nil {
				log.Fatalln(err)
			}
			targets := []string{}
			for _, target := range qr.Targets {
				targets = append(targets, target.Target)
			}
			queries := s.GetGrafanaQuery(qr.Range.From, qr.Range.To, targets)
			bRes, err := json.Marshal(queries)
			if err != nil {
				log.Fatalln(err)
			}
			w.Write(bRes)
			w.Header().Add("Content-Type", "application/json")
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
