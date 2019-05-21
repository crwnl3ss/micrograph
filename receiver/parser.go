package receiver

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/crwnl3ss/micrograph/storage"
)

var (
	metaDelimiter   = []byte(" ")
	spacesDelimiter = []byte(".")
)

type parsedPackage struct {
	rawNamespace string
	ts           int64
	metric       float64
}

func parseUDPRequest(b []byte) (string, *storage.DataPoint, error) {
	bs := bytes.Split(b, metaDelimiter)
	log.Println(string(b))
	if len(bs) < 3 {
		return nil, fmt.Errorf("invalid arguments count")
	}

	metric, err := strconv.ParseFloat(string(bs[1]), 10)
	if err != nil {
		return nil, fmt.Errorf("could not parse metric")
	}

	ts, err := strconv.ParseInt(string(bs[2]), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse timestamp")
	}
	return
	return &parsedPackage{rawNamespace: string(bs[0]), metric: metric, ts: ts}, nil
}
