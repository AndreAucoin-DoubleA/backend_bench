package server

import (
	"backend_bench/internal/model"
	"backend_bench/internal/routes"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func StartServer(config string, repo *model.UserRepository, wikiRepo *model.WikiRepository, jwtSecret string) {
	mux := http.NewServeMux()

	routes.RegisterRoutes(mux, repo, wikiRepo, jwtSecret)

	server := &http.Server{
		Addr:         ":" + config,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("Server running on :%s\n", config)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %s\n", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown failed: %v\n", err)
	}

	log.Println("Server exited properly")
}
