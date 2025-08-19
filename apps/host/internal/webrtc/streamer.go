package webrtc

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	pionwebrtc "github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/h264reader"
)

type Streamer struct {
	pc          *pionwebrtc.PeerConnection
	videoTrack  *pionwebrtc.TrackLocalStaticRTP
	payloadType uint8

	readyOnce sync.Once
	readyCh   chan struct{}
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

	sender, err := pc.AddTrack(videoTrack)
	if err != nil {
		return nil, fmt.Errorf("add track: %w", err)
	}
	// Drain RTCP to keep sender unblocked
	go func() {
		buf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := sender.Read(buf); rtcpErr != nil {
				return
			}
		}
	}()

	pc.OnICEConnectionStateChange(func(state pionwebrtc.ICEConnectionState) {
		log.Printf("ICE state: %s", state.String())
	})

	return &Streamer{
		pc:          pc,
		videoTrack:  videoTrack,
		payloadType: 96,
		readyCh:     make(chan struct{}),
	}, nil
}

func (s *Streamer) StartStream(stream io.ReadCloser, fps int) {
	go s.pumpStream(stream, fps)
}

func (s *Streamer) pumpStream(stream io.ReadCloser, fps int) {
	defer stream.Close()

	<-s.readyCh

	reader, err := h264reader.NewReader(stream)
	if err != nil {
		log.Printf("h264 reader: %v", err)
		return
	}

	pay := &codecs.H264Payloader{}
	seq := rtp.NewRandomSequencer()
	// 1200 bytes keeps us under typical 1500 MTU with headers
	pktizer := rtp.NewPacketizer(1200, s.payloadType, 0, pay, seq, 90000)

	// If you know real FPS, set frameDuration accordingly
	frameDuration := time.Second / time.Duration(fps)
	tsStep := uint32(90000 / uint32(fps))
	var ts uint32

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
				s.payloadType = uint8(c.PayloadType)
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
