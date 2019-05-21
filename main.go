package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/crwnl3ss/micrograph/receiver"
	"github.com/crwnl3ss/micrograph/web"
)

var s = receiver.NewStorage("hashmap")

func main() {
	var laddr = "127.0.0.1:8000"
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		cancel()
	}()

	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		err := receiver.Listen(ctx, laddr, s)
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		err := web.Listen(ctx, laddr, s)
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()
	log.Println("wait...")
	wg.Wait()
}
