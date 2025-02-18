package types

import "encoding/json"

// Client Messages
type Message struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

type ChatMessage struct {
	UnityID string `json:"unity_id"`
	Text    string `json:"text"`
	NpcId   string `json:"npcId"`
}

type EventMessage struct {
	UnityID      string `json:"unity_id"`
	EventType    string `json:"event_type"`
	EventDetails string `json:"event_details"`
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

type EventResponse struct {
	EventType string `json:"event_type"`
}
