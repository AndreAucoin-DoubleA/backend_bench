package config

import (
	"backend_bench/internal/model"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func GetConfig() model.Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file could not be loaded in local environment")
	}

	port := os.Getenv("PORT")
	stream := os.Getenv("STREAM_URL")
	cassandraHost := os.Getenv("CASSANDRA_HOST")
	cassandraPortStr := os.Getenv("CASSANDRA_PORT")
	cassandraPort, err := strconv.Atoi(cassandraPortStr)
	keyspaceKey := os.Getenv("KEYSPACE")
	jwtSecret := os.Getenv("JWT_SECRET")
	testKeyspace := os.Getenv("TEST_KEYSPACE")
	initialUserPassword := os.Getenv("INITIAL_USER_PASSWORD")
	initialEmail := os.Getenv("INITIAL_USER_EMAIL")

	if port == "" || stream == "" || cassandraHost == "" || cassandraPortStr == "" || keyspaceKey == "" || jwtSecret == "" || initialUserPassword == "" || initialEmail == "" {
		log.Fatal("One or more required environment variables are missing")
	}

	if err != nil {
		log.Fatalf("Invalid CASSANDRA_PORT: %v", err)
	}

	return model.Config{
		Port:          port,
		Stream:        stream,
		CassandraPort: cassandraPort,
		CassandraHost: cassandraHost,
		KeyspaceKey:   keyspaceKey,
		TestKeyspace:  testKeyspace,
		JWTSecret:     jwtSecret,
		Email:         initialEmail,
		Password:      initialUserPassword,
	}
}
