package udt

import (
	"log"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	start := time.Now()
//	for time.Now().Sub(start) < 10000 {
//	}
	time.Sleep(10 * time.Microsecond)
	end := time.Now()
	log.Printf("\nStart:   %s\nEnd:     %s\nElapsed:", start, end, end.Sub(start))
	
	id := 820208164
	log.Printf("%x", id)
	
	log.Printf("%d", 0x20c55bf6)
	
	a := 0x0100007f
	log.Printf("%d.%d.%d.%d", byte(a), byte(a>>8), byte(a>>16), byte(a>>24))
}
