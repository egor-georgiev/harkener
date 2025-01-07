package capture

import (
	"harkener/internal"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const snaplen = 1600

func Capture(interfaceName string, ignorePorts map[layers.TCPPort]struct{}, output chan layers.TCPPort, state *internal.State) {
	handle, err := pcap.OpenLive(interfaceName, snaplen, true, pcap.BlockForever)
	if err != nil {
		state.Errors <- err
		return
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		select {
		case <-state.Ctx.Done():
			return
		default:
			transportLayer := packet.TransportLayer()
			if transportLayer == nil {
				continue
			}

			tcp, ok := transportLayer.(*layers.TCP) // cast to TCP
			if !ok {
				continue
			}

			if _, exists := ignorePorts[tcp.DstPort]; exists {
				continue
			}

			if tcp.SYN && !tcp.ACK {
				output <- tcp.DstPort
			}

		}
	}
}
