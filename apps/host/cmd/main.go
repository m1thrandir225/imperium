package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/m1thrandir225/imperium/apps/host/internal/util"
	"github.com/m1thrandir225/imperium/apps/host/internal/video"
	"github.com/m1thrandir225/imperium/apps/host/internal/webrtc"
)

func main() {
	cfg, err := util.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	recorder := video.NewRecorder(&cfg.VideoConfig)
	streamer, err := webrtc.NewStreamer()
	if err != nil {
		log.Fatal(err)
	}

	signalingAddr := "127.0.0.1:8088"
	go func() {
		if err := webrtc.StartHTTPServer(streamer, signalingAddr); err != nil {
			log.Printf("Error starting signaling server: %v", err)
		}
	}()

	go func() {
		h264Stream, err := recorder.RecordScreen(nil)
		if err != nil {
			log.Printf("Error recording screen: %v", err)
		}
		streamer.StartStream(h264Stream, cfg.VideoConfig.FPS)
	}()

	//
	//uiManager := ui.NewUIManager(cfg)
	//
	//uiManager.RunUI()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	_ = recorder.StopRecording()
	_ = streamer.Close()
}
