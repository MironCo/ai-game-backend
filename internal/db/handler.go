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
