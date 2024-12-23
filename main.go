package main

import (
	"flag"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const snaplen = 1600

var interfaceName = flag.String("i", "eth0", "interface to listen on")

func listen(ifName string, output chan layers.TCPPort) {
	handle, err := pcap.OpenLive(ifName, snaplen, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		transportLayer := packet.TransportLayer()
		if transportLayer == nil {
			continue
		}

		tcp, ok := transportLayer.(*layers.TCP) // cast to TCP
		if !ok {
			continue
		}

		if tcp.SYN && !tcp.ACK {
			output <- tcp.DstPort
		}
	}
}

func main() {
	flag.Parse()
	portInfo := make(chan layers.TCPPort)
	go listen(*interfaceName, portInfo)
	for {
		port := <-portInfo
		log.Printf("%d", port)
	}
}
