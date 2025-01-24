package types

type Player struct {
	ID      string `json:"_id,omitempty" db:"id"`
	UnityID string `json:"unity_id" db:"unity_id"`
}

type DBChatMessage struct {
	ID          string `json:"_id,omitempty" db:"id"`
	UnityID     string `json:"unity_id" db:"unity_id"`
	MessageText string `json:"message" db:"message"`
	Sender      string `json:"sender" db:"sender"`
	SentTo      string `json:"sent_to" db:"sent_to"`
	CreatedAt   string `json:"created_at" db:"created_at"`
}
