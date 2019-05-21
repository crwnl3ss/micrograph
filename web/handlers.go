package web

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/crwnl3ss/micrograph/receiver"
)

type SearchRequest struct {
	Target string `json:"target,omitempty"`
}

type SearchResponse []string

// grafana datasorce index handler
func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(string("minigraph grafana exporter v0.1")))
}

// search returns list of xualified endpoint names (full path to the right
//most namespace)
// TODO: build b-tree of namespaces and returns them.
func search(s *receiver.HashmapStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sReq := &SearchRequest{}
		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(b, sReq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println(sReq.Target)
		var sRes SearchResponse = s.GetGrafanaTarggets()
		br, err := json.Marshal(sRes)
		if err != nil {
			w.Write(byte(err))
			w.WriteHeader(http.StatusBadGateway)
		}
	}
}
