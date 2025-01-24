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

func NewHandler() (*DBHandler, error) {
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
        INSERT INTO players (unity_id)
        VALUES ($1)
    `, req.UnityID)

	if err != nil {
		return fmt.Errorf("failed to create player: %w", err)
	}

	fmt.Printf("Inserted Player with Unity ID: %s\n", req.UnityID)
	return nil
}

func (h *DBHandler) GetPlayerByUnityId(unityID string) (*types.Player, error) {
	var player types.Player

	err := h.db.QueryRow(`
        SELECT id, unity_id 
        FROM players 
        WHERE unity_id = $1
    `, unityID).Scan(&player.ID, &player.UnityID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("player not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &player, nil
}
