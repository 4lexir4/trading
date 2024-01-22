package socket

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	lock  sync.RWMutex
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/bestspreads", s.handleBestSpreads)
	return http.ListenAndServe(":3000", nil)
}

func (s *Server) unregisterConn(ws *websocket.Conn) {
	s.lock.Lock()
	delete(s.conns, ws)
	s.lock.Unlock()

	fmt.Printf("unregister connection %s\n", ws.RemoteAddr())

	ws.Close()
}

func (s *Server) registerConn(ws *websocket.Conn) {
	s.lock.Lock()
	s.conns[ws] = true
	s.lock.Unlock()

	fmt.Printf("register connection %s\n", ws.RemoteAddr())

}

type X struct {
	Val int
}

func (s *Server) readLoop(ws *websocket.Conn) {
	defer func() {
		s.unregisterConn(ws)
	}()

	i := 0
	for {
		if err := ws.WriteJSON(X{Val: i}); err != nil {
			fmt.Println("write error:", err)
		}
		i++
		time.Sleep(time.Second)
	}
}

func (s *Server) handleBestSpreads(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panicln("websocket upgrade error:", err)
		return
	}

	s.registerConn(ws)

	go s.readLoop(ws)
}
