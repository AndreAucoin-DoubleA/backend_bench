package model

type User struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	CreatedAt    string `json:"created_at"`
}
