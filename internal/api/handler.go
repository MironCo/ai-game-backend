package api

import (
	"rd-backend/internal/db"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	dbHandler *db.DBHandler
}

func NewHandler(dbHandler *db.DBHandler) *APIHandler {
	return &APIHandler{
		dbHandler: dbHandler,
	}
}

func (h *APIHandler) HelloWorld(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World! This is a test (:",
	})
}
