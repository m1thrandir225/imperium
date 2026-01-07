package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/m1thrandir225/imperium/apps/host/internal/app"
	"github.com/m1thrandir225/imperium/apps/host/internal/ui"
)

const (
	AppName = "imperium-host"
)

func main() {
	_ = os.Setenv("LC_ALL", "C")
	_ = os.Setenv("FYNE_LANGUAGE", "en")

	application, err := app.New(AppName)
	if err != nil {
		log.Fatal(err)
	}

	application.WireAuthHandlers()
	application.WireSettingsHandlers()
	application.WireProgramsHandlers()
	application.WireHostHandlers()

	application.Start()

	uiManager := ui.NewManager(application.State, application.Bus)
	uiManager.RunUI()

	log.Println("UI closed")
	application.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

}
