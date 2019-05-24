package web

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/crwnl3ss/micrograph/storage"
)

// Listen accept incoming http requests
func Listen(ctx context.Context, laddr string, s *storage.HashmapStorage, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()

	srv := &http.Server{Addr: laddr, Handler: nil}
	log.Printf("listen http on: %s", laddr)
	go func() {
		<-ctx.Done()
		srv.Close()
	}()

	http.HandleFunc("/", index)
	http.HandleFunc("/search", search(s))
	http.HandleFunc("/query", query(s))

	return srv.ListenAndServe()
}
