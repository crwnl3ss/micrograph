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

// Listen ...
func Listen(ctx context.Context, laddr string, s *storage.HashmapStorage, wg *sync.WaitGroup) error {
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
	go func() {
		log.Printf("listen udp on: ", laddr)
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
				if bytes.Contains(buf, []byte("\n")) {
					n -= len([]byte("\n"))
				}
				t, dp, err := parseUDPRequest(buf[:n])
				log.Printf("from: %s size: %d body: %s", addr, n, buf[:n])
				if err != nil {
					log.Println(err)
					return
				}
				if err := s.InsertDataPoint(t, dp); err != nil {
					log.Printf("Could not insert datapoint into storage. Reason: %s", err)
					return
				}
			}()
		}
	}()
	return nil
}
