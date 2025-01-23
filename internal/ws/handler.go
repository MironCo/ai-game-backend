package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rd-backend/internal/ai"
	"rd-backend/internal/types"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	upgrader  websocket.Upgrader
	aiHandler *ai.AIHandler
}

func NewHandler() *WSHandler {
	return &WSHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		aiHandler: ai.NewHandler(),
	}
}

func (h *WSHandler) Handle(c *gin.Context) {
	ws, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("Upgrade Error: %w", err)
		return
	}

	defer ws.Close()

	for {
		var msg types.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Read Error: %v", err)
			}
			break
		}

		response := h.handleMessage(msg)
		ws.WriteJSON(response)
	}
}

func (h *WSHandler) handleMessage(msg types.Message) types.WSResponse {
	switch msg.Type {
	case "chat":
		var chatMsg types.ChatMessage
		fmt.Println(string(msg.Content))
		if err := json.Unmarshal(msg.Content, &chatMsg); err != nil {
			log.Fatal("Error Parsing Message to Chat Message: %w", err)
			return createErrorMessage("Invalid Chat Message")
		}
		return h.handleChatMessage(&chatMsg)
	default:
		return createErrorMessage("Unknown Message Type")
	}
}

// "chat"
func (h *WSHandler) handleChatMessage(msg *types.ChatMessage) types.WSResponse {
	completion, err := h.aiHandler.GetChatCompletion(msg.Text)
	if err != nil {
		return createErrorMessage(err.Error())
	}

	response := types.ChatResponse{
		Completion: completion,
		NpcId:      msg.NpcId,
	}

	content, _ := json.Marshal(response)
	return types.WSResponse{
		Type:    "chat",
		Content: content,
	}
}

func createErrorMessage(msg string) types.WSResponse {
	content, _ := json.Marshal(map[string]string{
		"error": msg,
	})
	return types.WSResponse{
		Type:    "error",
		Content: content,
	}
}
