package httpserver

import "github.com/m1thrandir225/imperium/apps/host/internal/webrtc"

func (s *Server) routes() {
	//auth-server endpoints
	s.mux.HandleFunc("/api/session/start", s.handleStartSessionRequest)
	s.mux.HandleFunc("/api/session/end", s.handleStop)
	s.mux.HandleFunc("/api/session/status", s.handleStatus)
	s.mux.HandleFunc("/api/session/programs", s.handleGetPrograms)

	// client endpoints
	s.mux.HandleFunc("/ws", s.wsServer.HandleWebSocket)

	webrtc.RegisterSignalingHandlers(s.mux, s.sessionService.WebRTCStreamer)
}
