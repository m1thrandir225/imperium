package httpserver

import (
	"log"
	"net/http"

	"github.com/m1thrandir225/imperium/apps/host/internal/session"
	ws "github.com/m1thrandir225/imperium/apps/host/internal/websocket"
)

type Server struct {
	mux            *http.ServeMux
	sessionService *session.Service
	wsServer       *ws.WebsocketHandler
}

func NewServer(sessionService *session.Service) *Server {
	s := &Server{
		mux:            http.NewServeMux(),
		sessionService: sessionService,
		wsServer:       ws.NewWebsocketHandler(),
	}

	s.routes()

	return s
}

func (s *Server) Serve(addr string) error {
	log.Print("Starting HTTP server on ", addr)
	return http.ListenAndServe(addr, s.mux)
}
