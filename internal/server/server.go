package server

import (
	"context"
	"encoding/binary"
	"harkener/internal"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	endpoint       = "/ws"
	sendBufferSize = 2
	writeTimeout   = time.Second * 5
)

var upgrader = websocket.Upgrader{}

type hub struct {
	ctx    context.Context
	spokes map[*spoke]context.CancelFunc
	lock   sync.RWMutex
}

type spoke struct {
	ctx    context.Context
	conn  *websocket.Conn
	input chan uint16
}

func (h *hub) run(input chan uint16) {
	for{
		// select{
		// case:
		// case:
		// }
	}
}

func newHub() (*hub, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return &hub{
		ctx:    ctx,
		spokes: make(map[*spoke]struct{}),
	}, cancel
}

func (h *hub) openSpoke(conn *websocket.Conn) *spoke {
	spoke := &spoke{
		conn:  conn,
		input: make(chan uint16),
	}
	h.lock.Lock()
	h.spokes[spoke] = struct{}{}
	h.lock.Unlock()

	return spoke
}


func (h *hub) closeSpoke(s *spoke) {
	close(s.input)
	s.conn.Close()

	h.lock.Lock()
	delete(h.spokes, s)
	h.lock.Unlock()
}



func (h *hub) close(){
	h.lock.Lock()
	for spoke := range h.spokes{
		close(s.input)
		s.conn.Close()
	}
}

// TODO: think about necessity of logging connection errors
func handler(h *hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	spoke := h.openSpoke(conn)
	defer h.closeSpoke(spoke)

	buf := make([]byte, sendBufferSize)
	for {
		select {
		case <-spoke.ctx.Done():
			return
		case message := <-spoke.input:
			binary.BigEndian.PutUint16(buf, message)
			conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			err := conn.WriteMessage(websocket.BinaryMessage, buf)
			if err != nil {
				return
			}
		}
	}

}

func Serve(portInfo chan uint16, bindAddr string, mainState *internal.State) {
	hub, cancel := newHub()
	go hub.run()
	http.HandleFunc(
		endpoint,
		func(w http.ResponseWriter, r *http.Request) {
			handler(hub, w, r)
		},
	)
	<-
}
