package types

type CreateUserResponse struct {
	UnityID string `json:"id"`
	Message string `json:"message"`
}

type LoginUserResponse struct {
	UnityID string `json:"id"`
	Message string `json:"message"`
}

type RegisterPhoneNumberResponse struct {
	UnityID     string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
}
