package types

import "encoding/json"

// Client Messages
type Message struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

type ChatMessage struct {
	Text  string `json:"text"`
	NpcId string `json:"npcId"`
}

// Server Reponses
type WSResponse struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

type ChatResponse struct {
	Completion string `json:"completion"`
	NpcId      string `json:"npcId"`
}
