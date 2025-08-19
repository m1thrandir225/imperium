package webrtc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type sdpMsg struct {
	SDP string `json:"sdp"`
}

func StartHTTPServer(s *Streamer, addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, htmlPage) // reuse existing page from package
	})

	mux.HandleFunc("/offer", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req sdpMsg
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		answerSDP, err := s.HandleOffer(req.SDP)
		if err != nil {
			log.Printf("offer error: %v", err)
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(sdpMsg{SDP: answerSDP})
	})

	log.Printf("H264 signaling on http://%s", addr)
	return http.ListenAndServe(addr, mux)
}

const htmlPage = `<!doctype html>
<html>
<head>
	<meta charset="utf-8"/>
	<meta name="viewport" content="width=device-width, initial-scale=1"/>
	<title>Imperium Host Stream</title>
	<style>
		html, body { margin: 0; padding: 0; background: #111; color: #ddd; }
		#wrap { display: flex; flex-direction: column; gap: 8px; padding: 12px; }
		video { width: 100%%; height: auto; background: #000; }
		button { padding: 8px 12px; }
	</style>
</head>
<body>
	<div id="wrap">
		<button id="start">Start viewer</button>
		<div>Renders remote H264 stream from the host.</div>
		<video id="v" autoplay playsinline controls></video>
	</div>
	<script>
	const btn = document.getElementById('start');
	const video = document.getElementById('v');

	btn.onclick = async () => {
		btn.disabled = true;
		const pc = new RTCPeerConnection();
		pc.addTransceiver('video', { direction: 'recvonly' });
		pc.ontrack = (ev) => {
			video.srcObject = ev.streams[0];
		};

		const offer = await pc.createOffer();
		await pc.setLocalDescription(offer);

		const res = await fetch('/offer', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ sdp: offer.sdp }),
		});
		if (!res.ok) {
			alert('Offer failed');
			return;
		}
		const data = await res.json();
		await pc.setRemoteDescription({ type: 'answer', sdp: data.sdp });
	};
	</script>
</body>
</html>`
