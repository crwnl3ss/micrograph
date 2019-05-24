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

var udpladdr = "0.0.0.0:6667"
var httpladdr = "0.0.0.0:6666"

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

	err := receiver.Listen(ctx, udpladdr, s, wg)
	if err != nil {
		log.Fatalf("udp listner error: %s", err.Error())
	}

	err = web.Listen(ctx, httpladdr, s, wg)
	if err != nil {
		log.Println(err)
	}
	wg.Wait()
	log.Println("stopped")
}
