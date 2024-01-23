package socket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/4lexir4/trading/orderbook"
	"github.com/bufbuild/buf/private/pkg/tmp"
	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string   `json:"type"`
	Topic   string   `json:"topic"`
	Symbols []string `json:"symbols"`
}

type MessageSpreads struct {
	Symbol  string                 `json:"symbol"`
	Spreads []orderbook.BestSpread `json:"spreads"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	crossSpreadch chan map[string][]orderbook.CrossSpread
	lock          sync.RWMutex
	conns         map[string]map[*WSConn]bool
}

func NewServer(crossSpreadCh chan map[string][]orderbook.CrossSpread) *Server {
	s := &Server{
		crossSpreadch: crossSpreadCh,
		conns:         make(map[string]map[*WSConn]bool),
	}
	for _, symbol := range symbols {
		s.conns[symbol] = map[*WSConn]bool{}
	}
	return s
}

func (s *Server) Start() error {
	http.HandleFunc("/", s.handleWS)
	go s.writeLoop()
	return http.ListenAndServe(":4000", nil)
}

func (s *Server) unregisterConn(ws *WSConn) {
	s.lock.Lock()
	for _, symbol := range ws.Symbols {
		delete(s.conns[symbol], ws)
	}
	s.lock.Unlock()

	fmt.Printf("unregister connection %s\n", ws.RemoteAddr())

	ws.Close()
}

func (s *Server) registerConn(ws *WSConn) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, symbol := range ws.Symbols {
		s.conns[symbol][ws] = true
		fmt.Printf("register connection to symbol %s %s\n", symbol, ws.RemoteAddr())
	}
}

func (s *Server) writeLoop() {
	for data := range s.crossSpreadch {
		for symbol, spreads := range data {
			for ws := range s.conns[symbol] {
				if err := ws.WriteJSON(spreads); err != nil {
					fmt.Println("socket write error", err)
					s.unregisterConn(ws)
				}
			}
		}
	}
}

func (s *Server) readLoop(ws *websocket.Conn) {
	defer ws.Close()

	for {
		msg := Message{}
		if err := ws.ReadJSON(&msg); err != nil {
			fmt.Println("socket read error:", err)
			break
		}
		if err := s.handleSocketMessage(ws, msg); err != nil {
			fmt.Println("handle msg error:", err)
			break
		}
	}
}

funct (s *Server) handleSocketMessage(ws *websocket.Conn, msg Message) error {
  wsConn := &WSConn {
    Conn: ws, 
    Topic: msg.Topic,
    Symbols: msg.Symbols,
  }

  s.registerConn(wsConn)
}

func (s *Server) handleBestSpreads(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket upgrade error:", err)
		return
	}

	s.registerConn(ws)
}
