package db

import (
	"database/sql"
	"fmt"
	"log"
	"rd-backend/internal/types"
)

func (h *DBHandler) GetPlayerByUnityId(unityID string) (*types.Player, error) {
	var player types.Player

	err := h.db.QueryRow(`
        SELECT id, unity_id, phone_number
        FROM players 
        WHERE unity_id = $1
    `, unityID).Scan(&player.ID, &player.UnityID, &player.PhoneNumber)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("player not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &player, nil
}

func (h *DBHandler) GetPlayerByPhoneNumber(phoneNumber string) (*types.Player, error) {
	var player types.Player

	//fmt.Println(phoneNumber)
	err := h.db.QueryRow(`
	SELECT id, unity_id, phone_number 
	FROM players 
	WHERE phone_number = $1`, phoneNumber).Scan(&player.ID, &player.UnityID, &player.PhoneNumber)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("player not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &player, nil
}

func (h *DBHandler) GetLastMessagesFromDB(unityID string, numberBack int) ([]types.DBChatMessage, error) {
	rows, err := h.db.Query(`
    SELECT message, sender, sent_to, created_at 
    	FROM messages 
    	WHERE unity_id = $1 
    	ORDER BY created_at DESC 
    	LIMIT $2
	`, unityID, numberBack)

	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	defer rows.Close()

	var messages []types.DBChatMessage
	for rows.Next() {
		var msg types.DBChatMessage
		if err := rows.Scan(&msg.MessageText, &msg.Sender, &msg.SentTo, &msg.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (h *DBHandler) GetLastEventsFromDB(unityID string, numberBack int) ([]types.DBPlayerEvent, error) {
	rows, err := h.db.Query(`
		SELECT unity_id, event_type, event_details, created_at
		FROM events 
		WHERE unity_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2
	`, unityID, numberBack)
	if err != nil {
		log.Printf("could not get events from database: %s", err.Error())
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	defer rows.Close()

	var events []types.DBPlayerEvent
	for rows.Next() {
		var event types.DBPlayerEvent
		if err := rows.Scan(&event.UnityID, &event.EventType, &event.EventDetails, &event.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (h *DBHandler) GetLastTextsFromDB(unityID string, npcNumber string, numberBack int) ([]types.DBTextMessage, error) {
	rows, err := h.db.Query(`
		SELECT unity_id, message, sender_number, receiver_number, player_number, created_at 
		FROM texts 
		WHERE unity_id = $1 
		AND (sender_number = $2 OR receiver_number = $2)
		ORDER BY created_at DESC 
		LIMIT $3
	`, unityID, npcNumber, numberBack)

	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	defer rows.Close()

	var messages []types.DBTextMessage
	for rows.Next() {
		var msg types.DBTextMessage
		if err := rows.Scan(&msg.UnityID, &msg.MessageText, &msg.SenderNumber, &msg.ReceiverNumber, &msg.PlayerNumber, &msg.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
