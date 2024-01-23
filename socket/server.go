package socket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/4lexir4/trading/orderbook"
	"github.com/gorilla/websocket"
)

type MessageSpreads struct {
	Symbol  string                 `json:"symbol"`
	Spreads []orderbook.BestSpread `json:"spreads"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	bsch  chan orderbook.BestSpread
	lock  sync.RWMutex
	conns map[*websocket.Conn]bool
}

func NewServer(bsch chan orderbook.BestSpread) *Server {
	return &Server{
		bsch:  bsch,
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/bestspreads", s.handleBestSpreads)
	go s.writeLoop()
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

func (s *Server) writeLoop() {
	for data := range s.bsch {
		for ws := range s.conns {
			ws.WriteJSON(data)
		}
	}
}

func (s *Server) handleBestSpreads(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket upgrade error:", err)
		return
	}

	s.registerConn(ws)
}
