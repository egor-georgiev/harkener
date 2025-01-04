package server

import (
	"log"

	"github.com/google/gopacket/layers"
)

func Serve(portInfo chan layers.TCPPort, bindAddress string) {
	for {
		port := <-portInfo
		log.Printf("%d", port)
	}
}
