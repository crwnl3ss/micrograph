package web

import (
	"context"
	"log"
	"net/http"

	"github.com/crwnl3ss/micrograph/receiver"
)

// Listen accept incoming http requests
func Listen(ctx context.Context, laddr string, s *receiver.HashmapStorage) error {
	srv := &http.Server{Addr: laddr, Handler: nil}
	log.Printf("Listen incoming http requests on: %s", laddr)
	http.HandleFunc("/", index)
	go func() {
		<-ctx.Done()
		srv.Close()
	}()
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
