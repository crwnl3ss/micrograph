package web

import (
	"context"
	"log"
	"net/http"

	"github.com/crwnl3ss/micrograph/storage"
)

// Listen accept incoming http requests
func Listen(ctx context.Context, laddr string, s *storage.HashmapStorage) error {
	srv := &http.Server{Addr: laddr, Handler: nil}
	log.Printf("Listen incoming http requests on: %s", laddr)
	http.HandleFunc("/", index)
	http.HandleFunc("/search", search(s))
	http.HandleFunc("/query", query(s))
	go func() {
		<-ctx.Done()
		srv.Close()
	}()
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
