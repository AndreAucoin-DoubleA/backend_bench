package routes

import (
	"backend_bench/internal/handler/login"
	"backend_bench/internal/model"
	"net/http"
)

func RegisterLoginRoutes(mux *http.ServeMux, repo *model.UserRepository, jwtSecret string) {
	mux.HandleFunc("/login", login.LoginHandler(repo, jwtSecret))
}
