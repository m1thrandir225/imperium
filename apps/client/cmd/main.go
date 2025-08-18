package main

import (
	"fmt"

	"github.com/m1thrandir225/imperium/apps/client/config"
	"github.com/m1thrandir225/imperium/apps/client/internal"
)

func main() {
	config := config.Load()

	router := internal.SetupRouter(config)

	router.Run(":8081")
	fmt.Println("This is the client app.")
}
