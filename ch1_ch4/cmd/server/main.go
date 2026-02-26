package main

import (
	"backend_bench/internal/config"
	"backend_bench/internal/db"
	"backend_bench/internal/model"
	"backend_bench/internal/server"
	"backend_bench/internal/service/wikiconsumer"
	"fmt"
	"os"
)

func main() {
	config := config.GetConfig()
	session := db.ConnectToCassandra(config.CassandraHost, config.CassandraPort, config.KeyspaceKey, config.Email, config.Password)
	defer session.Close()

	if os.Getenv("RUN_CI") != "true" {
		go wikiconsumer.StartWikiConsumer(config.Stream, session)
		fmt.Printf("Server is running on port: %s\n", config.Port)
		server.StartServer(config.Port, &model.UserRepository{Session: session}, &model.WikiRepository{Session: session}, config.JWTSecret)
	} else {
		fmt.Println("RUN_CI=true, skipping server and consumer for tests")
	}
}
