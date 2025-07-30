package main

import (
	"fmt"
	"github.com/m1thrandir225/imperium/apps/host/internal/util"
	"github.com/m1thrandir225/imperium/apps/host/internal/video"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	cfg, err := util.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	//
	//uiManager := ui.NewUIManager(cfg)
	//
	//uiManager.RunUI()

	outputPath := "output.mp4"

	recorder := video.NewRecorder(&cfg.VideoConfig)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool)
	go func() {
		err := recorder.RecordScreen(outputPath)
		if err != nil {
			log.Printf("Error starting recording: %v", err)
			done <- true
			return
		}
	}()
	select {
	case <-sigChan:
		fmt.Println("\nReceived interrupt signal. Stopping recording...")
	case <-time.After(10 * time.Second):
		fmt.Println("Recording duration reached. Stopping...")
	case <-done:
		fmt.Println("Recording stopped unexpectedly")
		return
	}
	err = recorder.StopRecording()
	if err != nil {
		log.Fatal(err)
	}

	absPath, err := filepath.Abs(outputPath)
	if err != nil {
		log.Printf("Warning: Could not resolve absolute path: %v", err)
		absPath = outputPath
	}
	fmt.Printf("Recording saved to: %s\n", absPath)

}
