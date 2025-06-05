package main

import (
	"fmt"

	"github.com/m1thrandir225/imperium/apps/client/internal/server"
)

func main() {
	router := server.SetupRouter()

	router.Run(":8080")
	fmt.Println("This is the client app.")
}
