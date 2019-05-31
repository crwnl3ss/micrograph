package web

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/crwnl3ss/micrograph/storage"
)

var startTime = time.Now()

type indexPage struct {
	Uptime    string
	Version   string
	KeysCount int
}

func index(s storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err := template.ParseFiles("./web/index.html")
		if err != nil {
			log.Fatalln(err)
		}
		tpl.Execute(w, indexPage{
			Uptime:    time.Since(startTime).String(),
			Version:   "0.0.1",
			KeysCount: len(s.GetKeys()),
		})
	}
}

// SearchRequest ...
type SearchRequest struct {
	Target string `json:"target,omitempty"`
}

// SearchResponse ...
type SearchResponse []string

// search returns list of grafana targets
func search(s storage.Storage) func(http.ResponseWriter, *http.Request) {
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
		var sRes SearchResponse = s.GetKeys()
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
func query(s storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			bReq, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`could not read request body`))
				return
			}
			defer r.Body.Close()
			qr := &QueryRequest{}
			if err := json.Unmarshal(bReq, qr); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`invalid request, valid request example: '{"range": {"from": 1559122100, "to": 1559122424}, "targets": [{"target": "first_target"}]}'`))
				return
			}
			targets := []string{}
			for _, target := range qr.Targets {
				targets = append(targets, target.Target)
			}
			type QueryResponse struct {
				Target     string
				Datapoints [][]interface{}
			}
			queries := s.RangeQuery(qr.Range.From, qr.Range.To, targets)
			response := []*QueryResponse{}
			for i := range queries {
				subResponse := &QueryResponse{Target: queries[i].Target, Datapoints: [][]interface{}{}}
				for j := range queries[i].DataPoints {
					subResponse.Datapoints = append(subResponse.Datapoints, []interface{}{queries[i].DataPoints[j].Data, queries[i].DataPoints[j].TS * 1000}) // gtafana expect ts in milliseconds
				}
				response = append(response, subResponse)
			}
			bRes, err := json.Marshal(response)
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
