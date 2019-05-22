package web

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/crwnl3ss/micrograph/storage"
)

// SearchRequest ...
type SearchRequest struct {
	Target string `json:"target,omitempty"`
}

// SearchResponse ...
type SearchResponse []string

// grafana datasorce index handler
func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(string("minigraph grafana exporter v0.1")))
}

// search returns list of xualified endpoint names (full path to the right
//most namespace)
func search(s *storage.HashmapStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			b, err := json.Marshal(s.GetGrafanaTargets())
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
			}
			w.Write(b)
			return
		}
		if r.Method == "POST" {
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

func query(s *storage.HashmapStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		reqB, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		qr := &QueryRequest{}
		if err := json.Unmarshal(reqB, qr); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		targets := []string{}
		for _, target := range qr.Targets {
			targets = append(targets, target.Target)
		}
		queries := s.GetGrafanaQuery(qr.Range.From, qr.Range.To, targets)
		log.Println(queries)
		json.Marshal(queries)
		resB, err := json.Marshal(queries)
		if err != nil {
			log.Println(err)
		}
		w.Write(resB)
		w.Header().Add("Content-Type", "application/json")
	}
}
