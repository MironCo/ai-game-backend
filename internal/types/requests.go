package types

type RegisterPlayerRequest struct {
	UnityID     string `json:"unity_id" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type LoginPlayerRequest struct {
	UnityID string `json:"unity_id" binding:"required"`
}

type RegisterPhoneNumberRequest struct {
	UnityID     string `json:"unity_id" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}
