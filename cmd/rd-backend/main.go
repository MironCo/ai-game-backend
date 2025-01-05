package main

import (
	"fmt"
	"os"

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

	router.GET("/", helloWorld)

	fmt.Printf("Server Running On Port 8080\n")
	router.Run(":" + port)
}

func helloWorld(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World! This is a test (:",
	})
}
