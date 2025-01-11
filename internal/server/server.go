package server

import (
	"context"
	"encoding/binary"
	"fmt"
	"harkener/internal"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	endpoint              = "/ws"
	sendBufferSize        = 2
	writeTimeout          = time.Second * 5
	serverShutdownTimeout = time.Second * 5
)

// TODO: parametrize
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type hub struct {
	spokes map[*spoke]struct{}
	lock   sync.RWMutex
}

type spoke struct {
	name  string
	input chan uint16
}

func newHub() *hub {
	return &hub{
		spokes: make(map[*spoke]struct{}),
	}
}

func (h *hub) run(input chan uint16) {
	defer h.close()
	for {
		msg, ok := <-input
		if ok {
			h.lock.RLock()
			for s := range h.spokes {
				select {
				case s.input <- msg:
				default:
					log.Printf("dropping packet for spoke %v\n", s.name)
				}
			}
			h.lock.RUnlock()
		} else {
			log.Printf("input channel is closed, closing the hub\n")
			return
		}
	}
}

func (h *hub) openSpoke(name string) *spoke {
	s := &spoke{
		name:  name,
		input: make(chan uint16),
	}
	h.lock.Lock()
	h.spokes[s] = struct{}{}
	h.lock.Unlock()

	return s
}

func (h *hub) closeSpoke(s *spoke) {
	close(s.input)

	h.lock.Lock()
	delete(h.spokes, s)
	h.lock.Unlock()
}

func (h *hub) close() {
	h.lock.Lock()
	for s := range h.spokes {
		close(s.input)
		delete(h.spokes, s)
	}
}

func handler(h *hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed while upgrading the connection: %v\n", err)
		return
	}
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	log.Printf("opened a connection for %v\n", addr)

	spoke := h.openSpoke(addr)
	defer h.closeSpoke(spoke)
	buf := make([]byte, sendBufferSize)
	for {
		message, ok := <-spoke.input
		if ok {
			binary.BigEndian.PutUint16(buf, message)
			conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			err := conn.WriteMessage(websocket.BinaryMessage, buf)
			if err != nil {
				log.Printf("failed while writing to %v: %v\n", addr, err)
				return
			}
		} else {
			log.Printf("input channel is closed, stopping the handler for %v\n", addr)
			return
		}
	}
}

func Serve(portInfo chan uint16, bindAddr string, state *internal.State, tlsCertPath, tlsKeyPath string) {
	server := &http.Server{Addr: bindAddr, Handler: nil}
	hub := newHub()
	go hub.run(portInfo)

	http.HandleFunc(
		endpoint,
		func(w http.ResponseWriter, r *http.Request) {
			handler(hub, w, r)
		},
	)

	if tlsCertPath != "" && tlsKeyPath != "" {
		log.Printf("starting the server on wss://%v%v", bindAddr, endpoint)
		go func() {
			err := server.ListenAndServeTLS(tlsCertPath, tlsKeyPath) // always non-nil
			if err != http.ErrServerClosed {
				state.Errors <- fmt.Errorf("got error from ws server: %v", err)
			}
		}()
	} else if tlsCertPath == "" && tlsKeyPath == "" {
		log.Printf("starting the server on ws://%v%v", bindAddr, endpoint)
		go func() {
			err := server.ListenAndServe() // always non-nil
			if err != http.ErrServerClosed {
				state.Errors <- fmt.Errorf("got error from ws server: %v", err)
			}
		}()
	} else {
		state.Errors <- fmt.Errorf("got tls cert path: %v and tls key path: %v, but both must be either present or absent", tlsCertPath, tlsKeyPath)
	}

	// at the same time, portInfo will be closed and hub.run() will exit
	<-state.Ctx.Done()
	log.Printf("shutting down the server\n")
	shutdownContext, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		log.Printf("got error during the server shutdown: %v\n", err)
	}
}
