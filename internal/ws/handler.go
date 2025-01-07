package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	upgrader websocket.Upgrader
}

func NewHandler() *WSHandler {
	return &WSHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *WSHandler) Handle(c *gin.Context) {
	ws, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Upgrade Error: %w", err)
		return
	}

	defer ws.Close()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Read Error: %w", err)
			break
		}

		response := h.handleMessage(msg)
		ws.WriteJSON(response)
	}
}

func (h *WSHandler) handleMessage(msg Message) Message {
	switch msg.Type {
	case "chat":
		var chatMsg ChatMessage
		if err := json.Unmarshal(msg.Content, &chatMsg); err != nil {
			log.Println("Error Parsing Message to Chat Message: %w", err)
			return createErrorMessage("Invalid Chat Message")
		}
		return h.handleChatMessage(chatMsg)
	default:
		return createErrorMessage("Uknown Message Type")
	}

}

func (h *WSHandler) handleChatMessage(msg ChatMessage) Message {
	response := ChatMessage{
		Text:  "Testing the backend for now, responding to " + msg.Text,
		NpcId: msg.NpcId,
	}

	content, _ := json.Marshal(response)
	return Message{
		Type:    "chat",
		Content: content,
	}
}

func createErrorMessage(msg string) Message {
	content, _ := json.Marshal(map[string]string{
		"error": msg,
	})
	return Message{
		Type:    "error",
		Content: content,
	}
}
