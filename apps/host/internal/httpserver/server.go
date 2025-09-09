package httpserver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/m1thrandir225/imperium/apps/host/internal/session"
	"github.com/m1thrandir225/imperium/apps/host/internal/webrtc"
	ws "github.com/m1thrandir225/imperium/apps/host/internal/websocket"
)

type Server struct {
	mux            *http.ServeMux
	sessionService *session.SessionService
	wsServer       *ws.WebsocketHandler
}

func NewServer(sessionService *session.SessionService) *Server {
	s := &Server{
		mux:            http.NewServeMux(),
		sessionService: sessionService,
		wsServer:       ws.NewWebsocketHandler(),
	}

	s.routes()

	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/session/start", s.handleStartSession)
	s.mux.HandleFunc("/api/session/status", s.handleStatus)
	s.mux.HandleFunc("/api/session/end", s.handleStop)
	s.mux.HandleFunc("/api/session/programs", s.handleGetPrograms)
	s.mux.HandleFunc("/ws", s.wsServer.HandleWebSocket)

	webrtc.RegisterSignalingHandlers(s.mux, s.sessionService.WebRTCStreamer)
}

func (s *Server) Serve(addr string) error {
	log.Print("Starting HTTP server on ", addr)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) handleStartSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req StartSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	_, err := s.sessionService.StartSession(r.Context(), req.ProgramID, req.ClientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Register current session with WS server
	if cur := s.sessionService.GetCurrentSession(); cur != nil {
		s.wsServer.RegisterSession(cur.ID, s.sessionService)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(StartSessionResponse{Status: "ACTIVE"})
}

func (s *Server) handleStatus(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := "IDLE"
	if s.sessionService.GetCurrentSession() != nil {
		status = "ACTIVE"
	}
	_ = json.NewEncoder(w).Encode(StatusResponse{Status: status})
}

func (s *Server) handleStop(w http.ResponseWriter, _ *http.Request) {
	_ = s.sessionService.EndSession()
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleGetPrograms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Println("Getting programs")

	log.Printf("request context: %v", r.Context())
	log.Printf("request url: %v", r.URL)
	log.Printf("request method: %v", r.Method)
	log.Printf("request headers: %v", r.Header)
	log.Printf("request body: %v", r.Body)
	log.Printf("request context: %v", r.Context())
	log.Printf("request url: %v", r.URL)
	log.Printf("request method: %v", r.Method)
	log.Printf("request headers: %v", r.Header)
	log.Printf("request body: %v", r.Body)

	programs, err := s.sessionService.GetPrograms()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(programs)
}
