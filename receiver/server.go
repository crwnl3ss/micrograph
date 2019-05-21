package receiver

import (
	"bytes"
	"context"
	"log"
	"net"
	"strings"
)

func processRequest(b []byte, s *HashmapStorage) {

}

// Listen ...
func Listen(ctx context.Context, laddr string, s *HashmapStorage) error {
	log.Printf("Listen for incoming udp packages on %s", laddr)
	pc, err := net.ListenPacket("udp", laddr)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		pc.Close()
	}()
	buf := make([]byte, 1024)
	for {
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Println("udp package listener was closed")
				return nil
			}
			return err
		}
		log.Printf("Receive message from %s size %d", addr, n)
		go func() {
			if bytes.Contains(buf, []byte("\n")) {
				n -= len([]byte("\n"))
			}
			pp, err := parse(buf[:n])
			if err != nil {
				log.Println(err)
				return
			}
			if err := s.insert(pp); err != nil {
				log.Println(err)
				return
			}
		}()
	}
}
