package types

import "time"

type Player struct {
	ID          string `json:"_id,omitempty" db:"id"`
	UnityID     string `json:"unity_id" db:"unity_id"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
}

type NPC struct {
	ID          string   `json:"npc_id"`
	Name        string   `json:"name"`
	PhoneNumber string   `json:"phone_number"`
	Location    string   `json:"location"`
	Occupation  string   `json:"occupation"`
	Traits      []string `json:"traits"`
	Quirks      []string `json:"quirks"`
	Goals       string   `json:"goals"`
	Backstory   string   `json:"backstory"`
	SpeechStyle string   `json:"speech_style"`
}

type DBChatMessage struct {
	ID          string `json:"_id,omitempty" db:"id"`
	UnityID     string `json:"unity_id" db:"unity_id"`
	MessageText string `json:"message" db:"message"`
	Sender      string `json:"sender" db:"sender"`
	SentTo      string `json:"sent_to" db:"sent_to"`
	CreatedAt   string `json:"created_at" db:"created_at"`
}

type DBTextMessage struct {
	ID             string    `json:"_id,omitempty" db:"id"`
	UnityID        string    `json:"unity_id" db:"unity_id"`
	MessageText    string    `json:"message" db:"message"`
	SenderNumber   string    `json:"sender_number" db:"sender_number"`
	ReceiverNumber string    `json:"receiver_number" db:"receiver_number"`
	PlayerNumber   string    `json:"player_number" db:"receiver_number"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type DBPlayerEvent struct {
	ID           string    `json:"_id,omitempty" db:"id"`
	UnityID      string    `json:"unity_id" db:"unity_id"`
	EventType    string    `json:"event_type" db:"event_type"`
	EventDetails string    `json:"event_details" db:"event_details"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
