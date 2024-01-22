package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Server struct {
	conn map[*websocket.Conn]bool
}

func (s *Server) Start() error {
	http.HandleFunc("/bestspreads", s.handleBestSpreads)
	return http.ListenAndServe(":3000", nil)
}

func (s *Server) handleBestSpreads(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panicln("websocket upgrade error:", err)
		return
	}
	defer ws.Close()
}
