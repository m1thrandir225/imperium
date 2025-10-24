package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"

	"github.com/m1thrandir225/imperium/apps/client/config"
	"github.com/m1thrandir225/imperium/apps/client/internal/server"
)

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Printf("Error opening browser: %v", err)
		fmt.Printf("Please open your browser and navigate to: %s\n", url)

	}
}

func main() {
	cfg := config.Load()

	if !config.IsConfigured() {
		fmt.Println("Application not configured. Please complete setup in your browser.")
	}

	port := ":8081"
	appURL := fmt.Sprintf("http://localhost%s", port)

	router := server.SetupRouter(cfg)

	go func() {
		if err := router.Run(port); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	openBrowser(appURL)

	fmt.Printf("Imperium Client is running on %s\n", appURL)
	if !config.IsConfigured() {
		fmt.Println("⚠️  Configuration required - please complete setup in your browser")
	}
	fmt.Println("Press Ctrl+C to stop the server")

	select {}
}
