package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/m1thrandir225/imperium/apps/host/internal/config"
	"github.com/m1thrandir225/imperium/apps/host/internal/ui"
)

func main() {
	os.Setenv("LC_ALL", "C")
	os.Setenv("FYNE_LANGUAGE", "en")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	uiManager := ui.NewUIManager(cfg)
	uiManager.RunUI()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

}
