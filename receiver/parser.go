package receiver

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/crwnl3ss/micrograph/storage"
)

var (
	bSpace = []byte(" ")
	bDot   = []byte(".")
)

type parsedPackage struct {
	rawNamespace string
	ts           int64
	metric       float64
}

func parseUDPRequest(b []byte) (string, *storage.DataPoint, error) {
	bs := bytes.Split(b, bSpace)
	if len(bs) < 3 {
		return "", nil, fmt.Errorf("invalid arguments count")
	}

	metric, err := strconv.ParseFloat(string(bs[1]), 10)
	if err != nil {
		return "", nil, fmt.Errorf("could not parse metric")
	}
	// truncate all numbers aftert `.` in timestamp
	timestamp, err := strconv.ParseInt(string(bytes.Split(bs[2], bDot)[0]), 10, 64)
	if err != nil {
		return "", nil, fmt.Errorf("could not parse timestamp")
	}
	return string(bs[0]), &storage.DataPoint{TS: timestamp, Data: metric}, nil
}
