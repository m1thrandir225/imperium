// package httpserver provides a httpserver that accepts requests
package httpserver

import (
	"log"
	"net/http"

	"github.com/m1thrandir225/imperium/apps/host/internal/session"
)

// Server provides the httpserver instance along with an EventBus to publish
// events and the sessionService
type Server struct {
	mux            *http.ServeMux
	sessionService *session.Service
	eventBus       interface{ Publish(topic string, data any) }
}

// NewServer returns a new instance of the Server struct
func NewServer(
	sessionService *session.Service,
	eventBus interface{ Publish(topic string, data any) },
) (*Server, error) {
	if sessionService == nil {
		return nil, InvalidSessionService
	}

	if eventBus == nil {
		return nil, InvalidEventBus
	}

	s := &Server{
		mux:            http.NewServeMux(),
		sessionService: sessionService,
		eventBus:       eventBus,
	}

	s.routes()

	return s, nil
}

// Serve starts the current httpserver instance
func (s *Server) Serve(addr string) error {
	log.Print("Starting HTTP server on ", addr)
	return http.ListenAndServe(addr, s.mux)
}

// Stop stops the current httpserver instance
func (s *Server) Stop() {
	log.Println("Stopping HTTP server")

	if s.mux != nil {
		s.mux = nil
	}

	if s.sessionService != nil {
		s.sessionService = nil
	}
}
