package login

import (
	"backend_bench/internal/auth"
	"backend_bench/internal/model"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(repo *model.UserRepository, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse JSON request
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Lookup user in Cassandra
		user, err := repo.GetUserByEmail(req.Email)
		if err != nil {
			fmt.Println("Error fetching user:", req.Email)
			http.Error(w, "Unable to find user", http.StatusUnauthorized)
			return
		}

		// Verify password
		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
			http.Error(w, "Password unable to be verified", http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		token, err := auth.GenerateJWT(user.UserID, user.Email, jwtSecret)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Return token
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(LoginResponse{Token: token}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
