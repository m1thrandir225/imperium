package main

import (
	"github.com/m1thrandir225/imperium/apps/host/config"
	"github.com/m1thrandir225/imperium/apps/host/internal/ui"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	uiManager := ui.NewUIManager(cfg)

	uiManager.RunUI()
}
