package db

import (
	"database/sql"
	"fmt"
	"os"
	"rd-backend/internal/types"

	_ "github.com/lib/pq"
)

type DBHandler struct {
	db *sql.DB
}

func NewDBHandler() (*DBHandler, error) {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping Postgres: %w", err)
	}
	fmt.Printf("Postgres Connected!\n")

	return &DBHandler{
		db: db,
	}, nil
}

func (h *DBHandler) Disconnect() error {
	if h.db != nil {
		if err := h.db.Close(); err != nil {
			return fmt.Errorf("failed to disconnect from Postgres: %w", err)
		}
	}
	return nil
}

func (h *DBHandler) CreatePlayer(req *types.RegisterPlayerRequest) error {
	_, err := h.db.Exec(`
        INSERT INTO players (unity_id, phone_number)
        VALUES ($1, $2)
    `, req.UnityID, req.PhoneNumber)

	if err != nil {
		return fmt.Errorf("failed to create player: %w", err)
	}

	fmt.Printf("Inserted Player with Unity ID: %s\n", req.UnityID)
	return nil
}

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

func (h *DBHandler) SetPlayerPhoneNumber(unityID string, phoneNumber string) (*types.Player, error) {
	var player types.Player

	err := h.db.QueryRow(`
    	UPDATE players 
    	SET phone_number = $1 
    	WHERE unity_id = $2 
    	RETURNING id, unity_id, phone_number`,
		phoneNumber, unityID).Scan(&player.ID, &player.UnityID, &player.PhoneNumber)

	if err != nil {
		return nil, fmt.Errorf("could not update player's phone number")
	}

	return &player, nil
}

func (h *DBHandler) AddMessageToDatabase(unityID string, messageText string, sender string, sentTo string) error {
	_, err := h.db.Exec(`
        INSERT INTO messages (unity_id, message, sender, sent_to)
        VALUES ($1, $2, $3, $4)
    `, unityID, messageText, sender, sentTo)

	if err != nil {
		fmt.Println("Error adding message!" + err.Error())
		return fmt.Errorf("failed to add message: %w", err)
	}

	//fmt.Printf("Inserted Player with Unity ID: %s\n", UnityID)
	return nil
}

func (h *DBHandler) AddTextToDatabase(unityID string, messageText string, senderNumber string, receiverNumber string, playerNumber string) error {
	_, err := h.db.Exec(`
        INSERT INTO texts (unity_id, message, sender_number, receiver_number, player_number)
        VALUES ($1, $2, $3, $4, $5)
    `, unityID, messageText, senderNumber, receiverNumber, playerNumber)

	if err != nil {
		fmt.Println("Error adding text message: " + err.Error())
		return fmt.Errorf("failed to add message: %w", err)
	}

	return nil
}

func (h *DBHandler) AddEventToDatabase(unityID string, eventType string, eventDetails string) error {
	_, err := h.db.Exec(`
		INSERT INTO events (unity_id, event_type, event_details)
		VALUES ($1, $2, $3)
	`, unityID, eventType, eventDetails)

	if err != nil {
		return fmt.Errorf("could not add event into database: %w", err)
	}

	return nil
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
