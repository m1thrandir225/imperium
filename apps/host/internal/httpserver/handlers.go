package httpserver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/m1thrandir225/imperium/apps/host/internal/session"
)

// handleStartSessionRequest represents the http handler for starting a session
func (s *Server) handleStartSessionRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req AuthServerSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	log.Printf("Received start session request from auth-server: %v", req)

	if s.sessionService.GetCurrentSession() != nil {
		response := AuthServerSessionResponse{
			Success: false,
			Error:   "Host already has an active session",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	session, err := s.sessionService.StartSession(r.Context(), session.StartSessionCommand{
		ProgramID:    req.ProgramID,
		ClientID:     req.ClientID,
		ClientName:   req.ClientName,
		Status:       req.Status,
		SessionID:    req.SessionID,
		HostID:       req.HostID,
		HostName:     req.HostName,
		SessionToken: req.SessionToken,
		WebrtcOffer:  req.WebrtcOffer,
		StartedAt:    req.StartedAt,
		CreatedAt:    req.CreatedAt,
	})
	if err != nil {
		log.Printf("Failed to start session: %v", err)
		response := AuthServerSessionResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	if s.eventBus != nil {
		s.eventBus.Publish("session.started", map[string]string{
			"SessionID":   session.ID,
			"ProgramName": session.WindowTitle,
			"ClientName":  session.ClientName,
		})
	}

	webrtcAnswer, err := s.sessionService.GenerateWebRTCAnswer(req.WebrtcOffer)
	if err != nil {
		log.Printf("Failed to generate WebRTC offer: %v", err)
		response := AuthServerSessionResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	response := AuthServerSessionResponse{
		Success:      true,
		WebrtcAnswer: webrtcAnswer,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
	return
}

// handleStatus represents the http handler for returning the current status
func (s *Server) handleStatus(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := "IDLE"
	if s.sessionService.GetCurrentSession() != nil {
		status = "ACTIVE"
	}
	_ = json.NewEncoder(w).Encode(StatusResponse{Status: status})
}

// handleStop represents the http handler for stopping the current session
func (s *Server) handleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := s.sessionService.EndSession()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if s.eventBus != nil {
		s.eventBus.Publish("session.ended", map[string]string{
			"SessionID": "ended",
		})
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleGetPrograms represents the http handler for returning the current host
// programs
func (s *Server) handleGetPrograms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	programs, err := s.sessionService.GetPrograms()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(programs)
}
