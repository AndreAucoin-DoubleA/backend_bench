package main

import (
	"backend_bench/internal/config"
	"backend_bench/internal/db"
	"backend_bench/internal/model"
	"backend_bench/internal/redpanda"
	"backend_bench/internal/server"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "wiki-changes"
	group := "wiki-consumer-group"
	consumerSession, err := redpanda.NewConsumerWithConfig(brokers, topic, group)
	if err != nil {
		log.Fatal(err)
	}
	config := config.GetConfig()
	session := db.ConnectToCassandra(config.CassandraHost, config.CassandraPort, config.KeyspaceKey, config.Email, config.Password)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Consuming from REDpanda")

	var closeConsumerOnce sync.Once
	closeConsumer := func() {
		closeConsumerOnce.Do(func() {
			consumerSession.Close()
		})
	}

	defer session.Close()
	defer closeConsumer()
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("interrupt received, shutting down consumer...")
		closeConsumer()
	}()

	if os.Getenv("RUN_CI") != "true" {
		go consumerSession.Consume(ctx, session)
		fmt.Printf("Server is running on port: %s\n", config.Port)
		server.StartServer(ctx, config.Port, &model.UserRepository{Session: session}, &model.WikiRepository{Session: session}, config.JWTSecret)
	} else {
		fmt.Println("RUN_CI=true, skipping server and consumer for tests")
	}
}
