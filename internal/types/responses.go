package types

type CreateUserResponse struct {
	UnityID string `json:"id"`
	Message string `json:"message"`
}

type LoginUserResponse struct {
	UnityID string `json:"id"`
	Message string `json:"message"`
}
