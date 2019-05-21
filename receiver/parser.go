package receiver

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
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

func parse(b []byte) (*parsedPackage, error) {
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
	log.Println(string(bs[2]), err)
	if err != nil {
		return nil, fmt.Errorf("could not parse timestamp")
	}
	return &parsedPackage{rawNamespace: string(bs[0]), metric: metric, ts: ts}, nil
}
