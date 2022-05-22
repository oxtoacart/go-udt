package main

import (
	"github.com/oxtoacart/go-udt/udt"
	"log"
	"time"
)

func main() {
	go server("localhost:47008")
	time.Sleep(200 * time.Millisecond)
	go client("localhost:47008")
	time.Sleep(50 * time.Second)
}

func server(addr string) {
	if _, err := udt.ListenUDT("udp", addr); err != nil {
		log.Fatalf("Unable to listen: %s", err)
	}
}

func client(addr string) {
	if _, err := udt.DialUDT("udp", nil, addr); err != nil {
		log.Fatalf("Unable to dial: %s", err)
	}
}
