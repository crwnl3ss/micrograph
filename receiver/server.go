package receiver

import (
	"bytes"
	"context"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/crwnl3ss/micrograph/storage"
)

// Listen udp packages on passed laddr, process with `parseUDPRequest` and
// save with `InsertDataPoint`
func Listen(ctx context.Context, laddr string, s storage.Storage, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()
	pc, err := net.ListenPacket("udp", laddr)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		pc.Close()
	}()
	buf := make([]byte, 1024)
	// listen incoming UDP packages in single goroutine
	go func() {
		log.Printf("listen udp on: %s", laddr)
		for {
			n, addr, err := pc.ReadFrom(buf)
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					log.Println("udp listner closed")
					return
				}
				return
			}
			go func() {
				for _, message := range bytes.Split(buf[:n], []byte("\n")) {
					// log.Printf("from: %s size: %d body: %s", addr, n, buf[:n])
					t, dp, err := parseUDPRequest(message)
					if err != nil {
						log.Printf("malformed message `%s` from %s, error: %s", buf[:n], addr, err)
						return
					}
					if err := s.InsertDataPoint(t, dp); err != nil {
						log.Printf("could not insert datapoint into storage %s", err)
						return
					}
				}
			}()
		}
	}()
	return nil
}
