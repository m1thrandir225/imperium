package webrtc

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/m1thrandir225/imperium/apps/host/internal/input"
	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	pionwebrtc "github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/h264reader"
)

type Streamer struct {
	pc               *pionwebrtc.PeerConnection
	videoTrack       *pionwebrtc.TrackLocalStaticRTP
	videoPayloadType uint8

	readyOnce  sync.Once
	readyCh    chan struct{}
	iceReadyCh chan struct{}
}

func NewStreamer() (*Streamer, error) {
	cfg := pionwebrtc.Configuration{
		ICEServers: []pionwebrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	}
	pc, err := pionwebrtc.NewPeerConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("create peer connection: %w", err)
	}

	videoTrack, err := pionwebrtc.NewTrackLocalStaticRTP(
		pionwebrtc.RTPCodecCapability{
			MimeType:    pionwebrtc.MimeTypeH264,
			ClockRate:   90000,
			SDPFmtpLine: "level-asymmetry-allowed=1;packetization-mode=1",
		},
		"video",
		"host",
	)
	if err != nil {
		return nil, fmt.Errorf("create track: %w", err)
	}

	videoSender, err := pc.AddTrack(videoTrack)
	if err != nil {
		return nil, fmt.Errorf("add track: %w", err)
	}

	// Drain RTCP to keep sender unblocked
	go func() {
		buf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := videoSender.Read(buf); rtcpErr != nil {
				return
			}
		}
	}()

	ordered := false
	maxRetrans := uint16(0)
	dataChannel, err := pc.CreateDataChannel("input", &pionwebrtc.DataChannelInit{
		Ordered:        &ordered,
		MaxRetransmits: &maxRetrans,
	})

	if err != nil {
		return nil, fmt.Errorf("create data channel: %w", err)
	}

	var msgCount uint64

	dataChannel.OnOpen(func() {
		log.Printf("input dc: open label=%q id=%d negotiated=%v readyState=%s",
			dataChannel.Label(), dataChannel.ID(), dataChannel.Negotiated(), dataChannel.ReadyState())
	})

	dataChannel.OnClose(func() {
		log.Printf("input dc: close label=%q", dataChannel.Label())
	})

	dataChannel.OnBufferedAmountLow(func() {
		log.Printf("input dc: bufferedAmountLow=%d", dataChannel.BufferedAmount())
	})

	dataChannel.OnMessage(func(msg pionwebrtc.DataChannelMessage) {
		msgCount++
		if msg.IsString {
			log.Printf("input dc: #%d wrong type=string len=%d (expect binary)", msgCount, len(msg.Data))
			return
		}

		if cmd, ok := input.DecodeInputCommand(msg.Data); ok {
			log.Printf("input dc: #%d decoded cmd=%+v", msgCount, cmd)
			input.HandleCommand(cmd)
		} else {
			log.Printf("Wrong message type. Expected binary.")
		}
	})

	streamer := &Streamer{
		pc:               pc,
		videoTrack:       videoTrack,
		videoPayloadType: 96,
		readyCh:          make(chan struct{}),
		iceReadyCh:       make(chan struct{}),
	}

	pc.OnICEConnectionStateChange((func(state pionwebrtc.ICEConnectionState) {
		log.Printf("ICE state: %s", state.String())

		// Signal when ICE connection is established
		if state == pionwebrtc.ICEConnectionStateConnected {
			log.Printf("ICE connection established, ready to stream video")
			close(streamer.iceReadyCh)
		} else if state == pionwebrtc.ICEConnectionStateFailed {
			log.Printf("ICE connection failed - may need TURN server")
		}
	}))

	pc.OnConnectionStateChange(func(state pionwebrtc.PeerConnectionState) {
		log.Printf("PeerConnection state: %s", state.String())
	})

	return streamer, nil
}

func (s *Streamer) StartStream(stream io.ReadCloser, fps int) {
	go s.pumpStream(stream, fps)
}

func (s *Streamer) pumpStream(stream io.ReadCloser, fps int) {
	defer stream.Close()

	log.Printf("Waiting for ice ready channel")
	<-s.iceReadyCh
	log.Printf("Ice ready channel closed")

	reader, err := h264reader.NewReader(stream)
	if err != nil {
		log.Printf("h264 reader: %v", err)
		return
	}

	pay := &codecs.H264Payloader{}
	seq := rtp.NewRandomSequencer()
	// 1200 bytes keeps us under typical 1500 MTU with headers
	pktizer := rtp.NewPacketizer(1200, s.videoPayloadType, 0, pay, seq, 90000)

	// If you know real FPS, set frameDuration accordingly
	frameDuration := time.Second / time.Duration(fps)
	tsStep := uint32(90000 / uint32(fps))
	var ts uint32

	log.Printf("Starting video stream at %d FPS", fps)

	for {
		nal, err := reader.NextNAL()
		if err != nil {
			if err != io.EOF {
				log.Printf("read NAL: %v", err)
			}
			return
		}

		// Packetize this NAL; payloader will fragment as FU-A when needed
		pkts := pktizer.Packetize(nal.Data, uint32(frameDuration))
		for _, p := range pkts {
			p.Timestamp = ts
			if err := s.videoTrack.WriteRTP(p); err != nil {
				log.Printf("WriteRTP: %v", err)
				return
			}
		}

		ts += tsStep
	}
}

// pumpAudioStream currently not used, as unable to get audio stream from the host.
// func (s *Streamer) pumpAudioStream(audioStream io.ReadCloser, sampleRate int) {
// 	defer audioStream.Close()

// 	log.Printf("Waiting for ice ready channel (audio)")
// 	<-s.iceReadyCh
// 	log.Printf("Ice ready channel closed (audio)")

// 	pay := &codecs.OpusPayloader{}
// 	seq := rtp.NewRandomSequencer()
// 	pktizer := rtp.NewPacketizer(1200, s.videoPayloadType, 0, pay, seq, uint32(sampleRate))

// 	frameDuration := 20 * time.Millisecond

// 	samplesPerFrame := uint32(sampleRate * 20 / 1000)
// 	log.Printf("Starting audio stream at %d Hz", sampleRate)

// 	buf := make([]byte, 1920)
// 	var ts uint32

// 	for {
// 		n, err := io.ReadFull(audioStream, buf)
// 		if err != nil {
// 			if err != io.EOF && err != io.ErrUnexpectedEOF {
// 				log.Printf("read audio: %v", err)
// 			}
// 			return
// 		}

// 		if n == 0 {
// 			continue
// 		}

// 		// Packetize the audio data
// 		pkts := pktizer.Packetize(buf[:n], uint32(frameDuration))
// 		for _, p := range pkts {
// 			p.Timestamp = ts
// 			if err := s.audioTrack.WriteRTP(p); err != nil {
// 				log.Printf("WriteRTP (audio): %v", err)
// 				return
// 			}
// 		}

// 		ts += samplesPerFrame
// 		time.Sleep(frameDuration) // Pace the audio
// 	}
// }

func (s *Streamer) HandleOffer(offerSDP string) (string, error) {
	offer := pionwebrtc.SessionDescription{Type: pionwebrtc.SDPTypeOffer, SDP: offerSDP}
	if err := s.pc.SetRemoteDescription(offer); err != nil {
		return "", fmt.Errorf("set remote: %w", err)
	}

	answer, err := s.pc.CreateAnswer(nil)
	if err != nil {
		return "", fmt.Errorf("create answer: %w", err)
	}

	gatherComplete := pionwebrtc.GatheringCompletePromise(s.pc)
	if err := s.pc.SetLocalDescription(answer); err != nil {
		return "", fmt.Errorf("set local: %w", err)
	}
	<-gatherComplete

	for _, sender := range s.pc.GetSenders() {
		for _, c := range sender.GetParameters().Codecs {
			if c.MimeType == pionwebrtc.MimeTypeH264 {
				s.videoPayloadType = uint8(c.PayloadType)
			}
		}
	}

	s.readyOnce.Do(func() { close(s.readyCh) })
	return s.pc.LocalDescription().SDP, nil
}

func (s *Streamer) Close() error {
	if s.pc != nil {
		return s.pc.Close()
	}
	return nil
}
