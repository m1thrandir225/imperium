package httpserver

import (
	"log"
	"net/http"

	"github.com/m1thrandir225/imperium/apps/host/internal/session"
)

type Server struct {
	mux            *http.ServeMux
	sessionService *session.Service
	eventBus       interface{ Publish(topic string, data any) }
}

func NewServer(sessionService *session.Service, eventBus interface{ Publish(topic string, data any) }) *Server {
	s := &Server{
		mux:            http.NewServeMux(),
		sessionService: sessionService,
		eventBus:       eventBus,
	}

	s.routes()

	return s
}

func (s *Server) Serve(addr string) error {
	log.Print("Starting HTTP server on ", addr)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) Stop() {
	log.Println("Stopping HTTP server")

	if s.mux != nil {
		s.mux = nil
	}

	if s.sessionService != nil {
		s.sessionService = nil
	}
}
