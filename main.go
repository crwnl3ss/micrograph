package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/crwnl3ss/micrograph/receiver"
	"github.com/crwnl3ss/micrograph/storage"
	"github.com/crwnl3ss/micrograph/web"
)

var udpladdr string
var httpladdr string

func init() {
	flag.StringVar(&udpladdr, "udpladdr", "0.0.0.0:6667", "--udpladdr=0.0.0.0:6667")
	flag.StringVar(&httpladdr, "httpladdr", "0.0.0.0:8000", "--httpladdr=0.0.0.0:8000")
	flag.Parse()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		cancel()
	}()

	wg := &sync.WaitGroup{}
	s := storage.GetStorageByType(ctx, "inmemory", wg)

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
