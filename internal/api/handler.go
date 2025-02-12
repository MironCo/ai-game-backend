package api

import (
	"net/http"
	"rd-backend/internal/db"
	"rd-backend/internal/types"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	dbHandler *db.DBHandler
}

func NewAPIHandler(dbHandler *db.DBHandler) *APIHandler {
	return &APIHandler{
		dbHandler: dbHandler,
	}
}

func (h *APIHandler) RegisterPlayer(c *gin.Context) {
	var req types.RegisterPlayerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.dbHandler.CreatePlayer(&req)

	c.JSON(http.StatusCreated, types.CreateUserResponse{
		UnityID: req.UnityID,
		Message: "Player Created Successfully",
	})
}

func (h *APIHandler) LoginPlayer(c *gin.Context) {
	var req types.LoginPlayerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	player, err := h.dbHandler.GetPlayerByUnityId(req.UnityID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.CreateUserResponse{
		UnityID: req.UnityID,
		Message: player.ID,
	})
}

func (h *APIHandler) RegisterPhoneNumber(c *gin.Context) {
	var req types.RegisterPhoneNumberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	player, err := h.dbHandler.SetPlayerPhoneNumber(req.UnityID, req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, types.RegisterPhoneNumberResponse{
		UnityID:     player.UnityID,
		PhoneNumber: player.PhoneNumber,
		Message:     "Phone number registered successfully!",
	})
}

func (h *APIHandler) HelloWorld(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World! This is a test (:",
	})
}

/*func (h *APIHandler) TestAIMessage(c *gin.Context) {
	var req types.ChatMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	completion, err := h.aiHandler.GetChatCompletion(req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "AI completion failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.ChatResponse{
		Completion: completion,
		NpcId:      req.NpcId,
	})
}*/
