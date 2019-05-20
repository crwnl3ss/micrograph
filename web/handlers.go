package web

import (
	"net/http"
)

// grafana datasorce index handler
func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(string("minigraph grafana exporter v0.1")))
}

// search returns list of xualified endpoint names (full path to the right
//most namespace)
// TODO: build b-tree of namespaces and returns them.
func search(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte("da"))
}
