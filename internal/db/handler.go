package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBHandler struct {
	client *mongo.Client
}

func NewHandler() (*DBHandler, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	options := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), options)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	fmt.Printf("MongoDB Connected!\n")

	return &DBHandler{
		client: client,
	}, err
}

func (h *DBHandler) Disconnect() error {
	if h.client != nil {
		if err := h.client.Disconnect(context.TODO()); err != nil {
			return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
		}
	}
	return nil
}

func (h *DBHandler) CreatePlayerDocument() {
	player := Player{
		Message: "Hello World!",
	}

	collection := h.client.Database("RdDatabase").Collection("Players")
	insertResult, err := collection.InsertOne(context.Background(), player)
	if err != nil {
		log.Fatal(err)
	}

	//Log
	fmt.Printf("Inserted Player with ID: %v\n", insertResult.InsertedID)
}
