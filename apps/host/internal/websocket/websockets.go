package websocket

import (
	"log"
	"net/http"
	"time"

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
				log.Printf("WebSocket origin check: %s", r.Header.Get("Origin"))
				return true // In production, implement proper origin checking
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		sessions: make(map[string]Session),
	}
}

func (s *WebsocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebSocket connection attempt from %s", r.RemoteAddr)
	log.Printf("WebSocket headers: %v", r.Header)

	// Check if this is a WebSocket upgrade request
	if r.Header.Get("Upgrade") != "websocket" {
		log.Printf("Not a WebSocket upgrade request, Upgrade header: %s", r.Header.Get("Upgrade"))
		http.Error(w, "Expected WebSocket upgrade", http.StatusBadRequest)
		return
	}

	// Check for required WebSocket headers
	if r.Header.Get("Connection") == "" {
		log.Printf("Missing Connection header")
		http.Error(w, "Missing Connection header", http.StatusBadRequest)
		return
	}

	log.Printf("WebSocket upgrade headers validated successfully")

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer func() {
		log.Printf("Closing WebSocket connection")
		conn.Close()
	}()

	log.Printf("WebSocket connection established successfully")

	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		log.Printf("No session ID provided in query parameters")
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInvalidFramePayloadData, "No session ID provided"))
		return
	}

	log.Printf("Session ID from query: %s", sessionID)

	sessionService, exists := s.sessions[sessionID]
	if !exists {
		log.Printf("Session not found: %s. Available sessions: %v", sessionID, s.getSessionIDs())
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInvalidFramePayloadData, "Session not found"))
		return
	}

	log.Printf("Session found, setting WebSocket connection")

	// Set the WebSocket connection in the session service
	sessionService.SetWebSocketConnection(conn)

	// Set up connection monitoring
	conn.SetPingHandler(func(message string) error {
		log.Printf("Received ping: %s", message)
		return conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second))
	})

	conn.SetPongHandler(func(message string) error {
		log.Printf("Received pong: %s", message)
		return nil
	})

	// Set read deadline to detect dead connections
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		log.Printf("Received pong, extending read deadline")
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	log.Printf("Starting WebSocket message loop for session: %s", sessionID)

	// Handle incoming messages
	messageCount := 0
	for {
		var msg InputMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket unexpected close error: %v", err)
			} else if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket normal close: %v", err)
			} else {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		messageCount++
		log.Printf("Received message #%d for session %s: %+v", messageCount, sessionID, msg)

		// Validate message
		if msg.SessionID != sessionID {
			log.Printf("Session ID mismatch: expected %s, got %s", sessionID, msg.SessionID)
			continue
		}

		// Process the input command
		log.Printf("Processing input command: %+v", msg.Command)
		sessionService.ProcessInputCommand(msg.Command)
	}

	log.Printf("WebSocket connection ended for session %s after %d messages", sessionID, messageCount)
}

func (s *WebsocketHandler) RegisterSession(sessionID string, sessionService Session) {
	log.Printf("Registering WebSocket session: %s", sessionID)
	s.sessions[sessionID] = sessionService
	log.Printf("Total registered sessions: %d", len(s.sessions))
}

func (s *WebsocketHandler) UnregisterSession(sessionID string) {
	log.Printf("Unregistering WebSocket session: %s", sessionID)
	delete(s.sessions, sessionID)
	log.Printf("Total registered sessions: %d", len(s.sessions))
}

func (s *WebsocketHandler) getSessionIDs() []string {
	ids := make([]string, 0, len(s.sessions))
	for id := range s.sessions {
		ids = append(ids, id)
	}
	return ids
}

// GetSessionCount returns the number of registered sessions
func (s *WebsocketHandler) GetSessionCount() int {
	return len(s.sessions)
}

// GetSessionIDs returns all registered session IDs
func (s *WebsocketHandler) GetSessionIDs() []string {
	return s.getSessionIDs()
}
