package main

import (
	"github.com/m1thrandir225/imperium/apps/host/internal/ui"
	"github.com/m1thrandir225/imperium/apps/host/internal/util"
	"log"
)

func main() {
	cfg, err := util.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	uiManager := ui.NewUIManager(cfg)

	uiManager.RunUI()
}
