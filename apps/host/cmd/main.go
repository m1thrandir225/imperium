package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/m1thrandir225/imperium/apps/host/internal/ui"
	"github.com/m1thrandir225/imperium/apps/host/internal/util"
)

func main() {
	cfg, err := util.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	uiManager := ui.NewUIManager(cfg)
	uiManager.RunUI()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

}
