package ws

import "encoding/json"

type Message struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

type ChatMessage struct {
	Text  string `json:"text"`
	NpcId string `json:"npcId"`
}
