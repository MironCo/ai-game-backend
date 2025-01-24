package main

import (
	"fmt"
	"log"
	"os"
	"rd-backend/internal/api"
	"rd-backend/internal/db"
	"rd-backend/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Set to release mode, but keep terminal logging
	gin.SetMode(gin.ReleaseMode)

	// Create router with default logger
	router := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default fallback
	}

	dbHandler, err := db.NewHandler()
	if err != nil {
		log.Fatal("Postgres Error: %w", err)
	}
	defer dbHandler.Disconnect()

	// Websockets
	wsHandler := ws.NewHandler()
	router.GET("/ws", wsHandler.Handle)

	apiHandler := api.NewHandler(dbHandler)
	router.GET("/hello", apiHandler.HelloWorld)
	router.POST("/register", apiHandler.RegisterPlayer)
	router.POST("/login", apiHandler.LoginPlayer)
	//router.POST("/test-ai", apiHandler.TestAIMessage)

	fmt.Println("Server Running On Port " + port)
	router.Run(":" + port)
}
