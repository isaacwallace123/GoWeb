package main

import (
	"github.com/isaacwallace123/GoWeb/core"
	"github.com/isaacwallace123/GoWeb/handlers"
	"log"
)

func main() {
	router := core.RegisterControllers(
		&handlers.UsersController{},
	)

	core.Use(core.LoggingMiddleware)

	log.Println("🚀 Server listening on http://localhost:8080")
	err := router.Listen(":8080")
	if err != nil {
		log.Fatal("❌ Server failed:", err)
	}
}
