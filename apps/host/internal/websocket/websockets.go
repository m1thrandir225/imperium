package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/m1thrandir225/imperium/apps/host/internal/input"
)

type Session interface {
	SetWebSocketConnection(conn *websocket.Conn)
	ProcessInputCommand(cmd input.InputCommand)
}

type WebsocketHandler struct {
	upgrader websocket.Upgrader
	sessions map[string]Session
}

type InputMessage struct {
	SessionID string             `json:"sessionId"`
	Command   input.InputCommand `json:"command"`
}

func NewWebsocketHandler() *WebsocketHandler {
	return &WebsocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // In production, implement proper origin checking
			},
		},
		sessions: make(map[string]Session),
	}
}

func (s *WebsocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		log.Printf("No session ID provided")
		return
	}

	sessionService, exists := s.sessions[sessionID]
	if !exists {
		log.Printf("Session not found: %s", sessionID)
		return
	}

	// Set the WebSocket connection in the session service
	sessionService.SetWebSocketConnection(conn)

	// Handle incoming messages
	for {
		var msg InputMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Process the input command
		sessionService.ProcessInputCommand(msg.Command)
	}
}

func (s *WebsocketHandler) RegisterSession(sessionID string, sessionService Session) {
	s.sessions[sessionID] = sessionService
}

func (s *WebsocketHandler) UnregisterSession(sessionID string) {
	delete(s.sessions, sessionID)
}
