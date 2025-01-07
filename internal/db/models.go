package db

type Player struct {
	ID      string `bson:"_id,omitempty"`
	Message string `bson:"message"`
}
