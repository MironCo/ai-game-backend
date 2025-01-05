package main

import (
	"fmt"
	"os"
	"rd-backend/internal/api"
	"rd-backend/internal/ws"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set to release mode, but keep terminal logging
	gin.SetMode(gin.ReleaseMode)

	// Create router with default logger
	router := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default fallback
	}

	// Websockets
	wsHandler := ws.NewHandler()
	router.GET("/ws", wsHandler.Handle)

	apiHandler := api.NewHandler()
	router.GET("/hello", apiHandler.HelloWorld)

	fmt.Printf("Server Running On Port 8080\n")
	router.Run(":" + port)
}
