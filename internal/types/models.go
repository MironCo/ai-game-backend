package types

type Player struct {
	ID      string `bson:"_id,omitempty"`
	UnityID string `bson:"unity_id"`
}
