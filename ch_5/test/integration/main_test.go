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
	// Load .env if exists
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Set defaults if CI didn't pass env
	host := os.Getenv("CASSANDRA_HOST")
	if host == "" {
		host = "localhost"
	}

	keyspace := os.Getenv("TEST_KEYSPACE")
	if keyspace == "" {
		log.Fatal("TEST_KEYSPACE must be set")
	}

	username := os.Getenv("CASSANDRA_USERNAME")
	if username == "" {
		username = "test@example.com"
	}

	password := os.Getenv("CASSANDRA_PASSWORD")
	if password == "" {
		password = "password"
	}

	// Connect to Cassandra
	testSession = db.ConnectToCassandra(host, 9042, keyspace, username, password)

	// Exit code after all tests run
	code := m.Run()

	// Cleanup
	if testSession != nil {
		testSession.Close()
	}

	os.Exit(code)
}
