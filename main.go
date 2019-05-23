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

var udpladdr = "127.0.0.1:6667"
var httpladdr = "127.0.0.1:6666"

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
		err := receiver.Listen(ctx, udpladdr, s)
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		err := web.Listen(ctx, httpladdr, s)
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()
	log.Println("ready to serve <3")
	wg.Wait()
}
