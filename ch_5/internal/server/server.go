package server

import (
	"backend_bench/internal/model"
	"backend_bench/internal/routes"
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

func StartServer(ctx context.Context, config string, repo *model.UserRepository, wikiRepo *model.WikiRepository, jwtSecret string) {
	mux := http.NewServeMux()

	routes.RegisterRoutes(mux, repo, wikiRepo, jwtSecret)

	server := &http.Server{
		Addr:         ":" + config,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	runServer := func() {
		log.Printf("Server running on :%s\n", config)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %s\n", err)
		}
	}

	if os.Getenv("RUN_CI") == "true" {
		go runServer()
		return
	}

	go runServer()

	<-ctx.Done()

	log.Println("Shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown failed: %v\n", err)
	}

	log.Println("Server exited properly")
}
