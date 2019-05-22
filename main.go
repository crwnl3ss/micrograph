package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/crwnl3ss/micrograph/receiver"
	"github.com/crwnl3ss/micrograph/storage"
	"github.com/crwnl3ss/micrograph/web"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		cancel()
	}()

	wg := &sync.WaitGroup{}
	s := storage.NewStorage(ctx, "inmemory", wg)
	go func() {
		wg.Add(1)
		err := receiver.Listen(ctx, "127.0.0.1:8000", s)
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		err := web.Listen(ctx, "127.0.0.1:8080", s)
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()
	log.Println("wait...")
	wg.Wait()
}
