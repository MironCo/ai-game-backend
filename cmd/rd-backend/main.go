package main

import (
	"fmt"
	"log"
	"os"
	"rd-backend/internal/ai"
	"rd-backend/internal/ai/npc"
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

	// NPC Config
	npcs, err := npc.LoadNPCConfig("internal/config/npc.json")
	if err != nil {
		log.Fatal("Cannot Load NPC Config: %w", err)
	}
	npcPhoneNumbers := npc.BuildPhoneIndex(npcs)

	// Database
	dbHandler, err := db.NewDBHandler()
	if err != nil {
		log.Fatal("Postgres Error: %w", err)
	}
	defer dbHandler.Disconnect()

	// AI
	aiHandler := ai.NewAIHandler(&npcs, &npcPhoneNumbers)

	// Websockets
	wsHandler := ws.NewWebsocketHandler(dbHandler, aiHandler)
	router.GET("/ws", wsHandler.Handle)

	//Texting TODO
	textingHandler := api.NewTextingHandler(dbHandler, aiHandler)
	//go textingHandler.SendSMSBasic()

	// API
	apiHandler := api.NewAPIHandler(dbHandler)
	router.GET("/hello", apiHandler.HelloWorld)
	router.POST("/register", apiHandler.RegisterPlayer)
	router.POST("/login", apiHandler.LoginPlayer)
	router.POST("/sms/receive", textingHandler.ReceiveSMS)
	//router.POST("/test-ai", apiHandler.TestAIMessage)

	fmt.Println("Server Running On Port " + port)
	router.Run(":" + port)
}
