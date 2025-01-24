package db

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"rd-backend/internal/types"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// type DBHandler struct {
// 	client *mongo.Client
// }

// func NewHandler() (*DBHandler, error) {
// 	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
// 	options := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)

// 	client, err := mongo.Connect(context.TODO(), options)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
// 	}

// 	// Ping to verify connection
// 	if err := client.Ping(context.TODO(), nil); err != nil {
// 		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
// 	}
// 	fmt.Printf("MongoDB Connected!\n")

// 	return &DBHandler{
// 		client: client,
// 	}, err
// }

// func (h *DBHandler) Disconnect() error {
// 	if h.client != nil {
// 		if err := h.client.Disconnect(context.TODO()); err != nil {
// 			return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
// 		}
// 	}
// 	return nil
// }

// func (h *DBHandler) CreatePlayerDocument(req *types.RegisterPlayerRequest) {
// 	player := types.Player{
// 		UnityID: req.UnityID,
// 	}

// 	collection := h.client.Database("RdDatabase").Collection("Players")
// 	insertResult, err := collection.InsertOne(context.Background(), player)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	//Log
// 	fmt.Printf("Inserted Player with ObjectID: %v\n", insertResult.InsertedID)
// }

// func (h *DBHandler) GetPlayerByUnityId(unityID string) (*types.Player, error) {
// 	collection := h.client.Database("RdDatabase").Collection("Players")

// 	filter := bson.M{"unity_id": unityID}

// 	var player types.Player
// 	err := collection.FindOne(context.TODO(), filter).Decode(&player)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, fmt.Errorf("player not found")
// 		}
// 		return nil, fmt.Errorf("database error: %v", err)
// 	}

// 	return &player, nil
// }
