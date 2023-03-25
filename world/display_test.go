package world

import (
	"strconv"
	"testing"
	"time"
)

func TestNewDisplayClient(t *testing.T) {
	dc := NewDisplayClient("grid")

	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		dc.Send([]byte(strconv.Itoa(i)))
	}
}
