package webrtc

import (
	"encoding/json"
	"net/http"
)

type sdpMsg struct {
	SDP string `json:"sdp"`
}

func RegisterSignalingHandlers(mux *http.ServeMux, getStreamer func() *Streamer) {
	mux.HandleFunc("/api/session/webrtc/offer", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		streamer := getStreamer()
		if streamer == nil {
			http.Error(w, "no active session", http.StatusServiceUnavailable)
			return
		}

		var req sdpMsg
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		answerSDP, err := streamer.HandleOffer(req.SDP)
		if err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(sdpMsg{SDP: answerSDP})
	})
}
