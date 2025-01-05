package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", helloWorld)
}

func helloWorld(*gin.Context) {
	fmt.Printf("Hello world!\n")
}
