package integration

import (
	"backend_bench/internal/db"
	"log"
	"os"
	"testing"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/joho/godotenv"
)

var testSession *gocql.Session

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	host := os.Getenv("CASSANDRA_HOST")
	if host == "" {
		host = "localhost"
	}

	keyspace := os.Getenv("TEST_KEYSPACE")
	if keyspace == "" {
		log.Fatal("TEST_KEYSPACE must be set")
	}

	testSession = db.ConnectToCassandra(host, 9042, keyspace, "test@example.com", "password")

	code := m.Run()

	testSession.Close()
	os.Exit(code)
}
