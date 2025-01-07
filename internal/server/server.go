package server

import (
	"encoding/binary"
	"fmt"
	"harkener/internal"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/gopacket/layers"
)

const (
	serverBufferSize         = 1024
	consecutiveErrsThreshold = 100
	writeTimeout             = time.Second * 1
	portBufferSize           = 2
)

type clientInfo struct {
	Lock sync.Mutex
	Info map[*net.UDPAddr]struct{}
}

func newClientInfo() *clientInfo {
	return &clientInfo{
		Info: make(map[*net.UDPAddr]struct{}),
	}
}

func acceptConnections(info *clientInfo, state *internal.State, conn *net.UDPConn) {

	var buf [serverBufferSize]byte
	var consecutiveErrsCount int
	for {
		select {
		case <-state.Ctx.Done():
			return
		default:
			if consecutiveErrsCount >= consecutiveErrsThreshold {
				state.Errors <- fmt.Errorf("consecutive errors threshold reached: %v", consecutiveErrsThreshold)
				return
			}

			_, addr, err := conn.ReadFromUDP(buf[0:])
			if err != nil {
				log.Printf("got error while reading from %v: %v\n", addr, err)
				consecutiveErrsCount++
				continue
			}
			consecutiveErrsCount = 0

			info.Lock.Lock()
			info.Info[addr] = struct{}{}
			info.Lock.Unlock()

		}
	}

}

func writeToConnections(portInfo chan layers.TCPPort, clientInfo *clientInfo, state *internal.State, conn *net.UDPConn) {
	buf := make([]byte, portBufferSize)
	for port := range portInfo {
		select {
		case <-state.Ctx.Done():
			return
		default:
			clientInfo.Lock.Lock()
			for addr := range clientInfo.Info {
				conn.SetWriteDeadline(time.Now().Add(writeTimeout))
				binary.BigEndian.PutUint16(buf, uint16(port))
				_, err := conn.WriteToUDP(buf, addr)
				// TODO: improved error handling and logging
				if err != nil {
					delete(clientInfo.Info, addr)
				}
			}
			clientInfo.Lock.Unlock()
		}
	}
}

func Serve(portInfo chan layers.TCPPort, bindAddress string, state *internal.State) {
	udpAddr, err := net.ResolveUDPAddr("udp", bindAddress)
	if err != nil {
		state.Errors <- fmt.Errorf("failed while binding to address %v: %v", udpAddr, err)
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		state.Errors <- fmt.Errorf("failed while creating a connection: %v\n", err)
		return
	}
	defer conn.Close()

	clientInfo := newClientInfo()

	go acceptConnections(clientInfo, state, conn)
	go writeToConnections(portInfo, clientInfo, state, conn)

	<-state.Ctx.Done()
}
