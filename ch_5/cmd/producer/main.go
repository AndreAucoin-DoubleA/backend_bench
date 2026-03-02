package main

import (
	"backend_bench/internal/model"
	"backend_bench/internal/redpanda"
	wikiconnection "backend_bench/internal/service/wikiConnection"
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	streamURL := "https://stream.wikimedia.org/v2/stream/recentchange"
	brokers := []string{"localhost:9092"}
	topic := "wiki-changes"

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	producer := model.NewProducer(brokers, topic)

	scanner, response := wikiconnection.ConnectToWikiStream(ctx, streamURL)
	if scanner == nil || response == nil {
		log.Fatal("Failed to connect to wiki stream")
	}

	var shutdownOnce sync.Once
	shutdown := func() {
		shutdownOnce.Do(func() {
			_ = response.Body.Close()
			producer.Close()
		})
	}
	defer shutdown()

	go func() {
		<-ctx.Done()
		fmt.Println("\nInterrupt received, shutting down producer...")
		shutdown()
	}()

	redpanda.ConnectRedpandaProducer(scanner, producer, ctx)
	fmt.Println("Producer stopped")
}
