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
	APP_NAME = "imperium-host"
)

func main() {
	os.Setenv("LC_ALL", "C")
	os.Setenv("FYNE_LANGUAGE", "en")

	application, err := app.New(APP_NAME)
	if err != nil {
		log.Fatal(err)
	}

	application.WireAuthHandlers()
	application.WireSettingsHandlers()
	application.WireProgramsHandlers()
	application.WireHostHandlers()

	application.Start()

	uiManager := ui.NewUIManager(application.State, application.Bus)
	uiManager.RunUI()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

}
