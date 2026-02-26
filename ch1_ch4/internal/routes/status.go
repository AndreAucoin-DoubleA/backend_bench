package routes

import (
	"backend_bench/internal/handler/status"
	"backend_bench/internal/middleware"
	"backend_bench/internal/model"

	"net/http"
)

func RegisterStatusRoutes(mux *http.ServeMux, jwtSecret string, wikiRepo *model.WikiRepository) {
	mux.Handle("/status", middleware.AuthMiddleware(http.HandlerFunc(status.StatusHandler(wikiRepo)), jwtSecret))
}
