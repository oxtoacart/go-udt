package main

import (
	"github.com/oxtoacart/go-udt/udt"
	"log"
	"net"
	"time"
)

func main() {
	if addr, err := net.ResolveUDPAddr("udp", "localhost:47008"); err != nil {
		log.Fatalf("Unable to resolve address: %s", err)
	} else {
		go server(addr)
		go client(addr)
		
		time.Sleep(5 * time.Second)
	}
}

func server(addr *net.UDPAddr) {
	if _, err := udt.ListenUDT("udp", addr); err != nil {
		log.Fatalf("Unable to listen: %s", err)
	}
}

func client(addr *net.UDPAddr) {
	if _, err := udt.DialUDT("udp", nil, addr); err != nil {
		log.Fatalf("Unable to dial: %s", err)
	}
}
