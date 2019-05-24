package receiver

import (
	"fmt"
	"testing"
	"time"
)

func TestParser(t *testing.T) {
	ts1 := time.Date(2019, 11, 16, 13, 25, 7, 0, time.UTC)
	t1 := "a.b.c.d"
	d1 := 0.1
	pckg1 := fmt.Sprintf("%s %f %d", t1, d1, ts1.Unix())
	tResp, dpResp, err := parseUDPRequest([]byte(pckg1))
	if err != nil {
		t.Errorf("unexpected non-nil error: %s", err)
	}
	if tResp != t1 {
		t.Errorf("targets mismatch: %s != %s", t1, tResp)
	}
	if dpResp.TS != ts1.Unix() || dpResp.Data != d1 {
		t.Errorf("malformed DataPoint created")
	}
}
