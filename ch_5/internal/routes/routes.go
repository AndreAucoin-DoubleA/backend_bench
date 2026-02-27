package routes

import (
	"backend_bench/internal/model"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, repo *model.UserRepository, wikiRepo *model.WikiRepository, jwtSecret string) {
	RegisterStatusRoutes(mux, jwtSecret, wikiRepo)
	RegisterLoginRoutes(mux, repo, jwtSecret)
}
