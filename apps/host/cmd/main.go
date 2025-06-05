package main

import (
	"github.com/m1thrandir225/imperium/apps/host/internal/ui"
)

func main() {
	app := ui.NewUIManager()

	app.RunUI()
}
