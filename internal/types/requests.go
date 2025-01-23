package types

type RegisterPlayerRequest struct {
	UnityID string `json:"unity_id" binding:"required"`
}

type LoginPlayerRequest struct {
	UnityID string `json:"unity_id" binding:"required"`
}
