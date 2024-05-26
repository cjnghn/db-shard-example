package main

import (
	"log"

	"github.com/cjnghn/db-shard-example/internal/config"
	"github.com/cjnghn/db-shard-example/internal/db"
	"github.com/cjnghn/db-shard-example/internal/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize configuration
	cfg := config.GetConfig()

	// Initialize shards
	err := db.InitShards(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize shards: %v", err)
	}

	// Create Echo instance
	e := echo.New()

	// Define routes
	e.POST("/users", handlers.CreateUserHandler)
	e.GET("/users/:id", handlers.GetUserHandler)
	e.GET("/users", handlers.GetAllUsersHandler)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
