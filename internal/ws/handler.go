package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"rd-backend/internal/ai"
	"rd-backend/internal/db"
	"rd-backend/internal/types"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	upgrader  websocket.Upgrader
	aiHandler *ai.AIHandler
	dbHandler *db.DBHandler
}

func NewWebsocketHandler(dbHandler *db.DBHandler, aiHandler *ai.AIHandler) *WSHandler {
	return &WSHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		aiHandler: aiHandler,
		dbHandler: dbHandler,
	}
}

func (h *WSHandler) Handle(c *gin.Context) {
	// Get and validate Unity ID
	unityID := c.Query("unity_id")
	exists, err := h.dbHandler.GetPlayerByUnityId(unityID)
	if err != nil || exists == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "player not found"})
		return
	}

	// If auth passes, upgrade to WebSocket
	ws, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		createErrorMessage("Upgrade Error: " + err.Error())
		return
	}

	defer ws.Close()

	for {
		var msg types.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				createErrorMessage("Upgrade Error: " + err.Error())
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
		//fmt.Println(string(msg.Content))
		if err := json.Unmarshal(msg.Content, &chatMsg); err != nil {
			log.Printf("Error Parsing Message to Chat Message: %v", err)
			return createErrorMessage("Invalid Chat Message")
		}
		return h.handleChatMessage(&chatMsg)
	case "system":
		var systemMsg types.ChatMessage
		if err := json.Unmarshal(msg.Content, &systemMsg); err != nil {
			log.Printf("Error Parsing Message to System Message: %v", err)
			return createErrorMessage("Invalid System Message")
		}
		return h.handleSystemMessage(&systemMsg)
	case "event":
		return createErrorMessage("Event... event (websocket event) not yet implemented")
	default:
		return createErrorMessage("Unknown Message Type")
	}
}

// "chat"
func (h *WSHandler) handleChatMessage(msg *types.ChatMessage) types.WSResponse {
	history, err := h.dbHandler.GetLastMessagesFromDB(msg.UnityID, 4)
	if err != nil {
		return createErrorMessage(err.Error())
	}

	h.dbHandler.AddMessageToDatabase(msg.UnityID, msg.Text, "player", msg.NpcId)

	completion, err := h.aiHandler.GetChatCompletion(msg.Text, history, "user", msg.NpcId)
	if err != nil || completion == nil {
		return createErrorMessage(err.Error())
	}

	response := types.ChatResponse{
		Completion: *completion,
		NpcId:      msg.NpcId,
	}

	h.dbHandler.AddMessageToDatabase(msg.UnityID, response.Completion, msg.NpcId, "player")

	content, _ := json.Marshal(response)

	return types.WSResponse{
		Type:    "chat",
		Content: content,
	}
}

// "system"
func (h *WSHandler) handleSystemMessage(msg *types.ChatMessage) types.WSResponse {
	history, err := h.dbHandler.GetLastMessagesFromDB(msg.UnityID, 4)
	if err != nil {
		return createErrorMessage("Could not get last messages from Database")
	}

	completion, err := h.aiHandler.GetChatCompletion(msg.Text, history, "system", msg.NpcId)
	if err != nil || completion == nil {
		return createErrorMessage(err.Error())
	}

	response := types.ChatResponse{
		Completion: *completion,
		NpcId:      msg.NpcId,
	}

	h.dbHandler.AddMessageToDatabase(msg.UnityID, response.Completion, msg.NpcId, "player")

	content, _ := json.Marshal(response)

	return types.WSResponse{
		Type:    "system",
		Content: content,
	}
}

/*func (h *WSHandler) handleEventMessage(msg *types.EventMessage) types.WSResponse {
	event := types.DBPlayerEvent{
		UnityID:   msg.UnityID,
		EventType: msg.EventType,

	}
}*/

func createErrorMessage(msg string) types.WSResponse {
	content, _ := json.Marshal(map[string]string{
		"error": msg,
	})
	return types.WSResponse{
		Type:    "error",
		Content: content,
	}
}
