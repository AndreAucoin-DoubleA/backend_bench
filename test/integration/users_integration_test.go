package integration

import (
	"backend_bench/internal/handler/login"
	"backend_bench/internal/model"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateAndGetUser(t *testing.T) {
	repo := &model.UserRepository{Session: testSession}

	// Clean table
	if err := testSession.Query("TRUNCATE user_by_email").Exec(); err != nil {
		t.Fatalf("Failed to truncate table: %v", err)
	}

	id := gocql.TimeUUID()
	createdAt := time.Now().Format(time.RFC3339)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if err := testSession.Query(`
        INSERT INTO user_by_email (email, id, password_hash, created_at)
        VALUES (?, ?, ?, ?)`,
		"test@example.com", id, hashedPassword, createdAt,
	).Exec(); err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	got, err := repo.GetUserByEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if got.Email != "test@example.com" {
		t.Fatalf("Expected %s, got %s", "test@example.com", got.Email)
	}
}

func TestLoginSuccess(t *testing.T) {
	repo := &model.UserRepository{Session: testSession}

	// Clean table
	if err := testSession.Query("TRUNCATE user_by_email").Exec(); err != nil {
		t.Fatalf("Failed to truncate table: %v", err)
	}

	id := gocql.TimeUUID()
	createdAt := time.Now().Format(time.RFC3339)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if err := testSession.Query(`
        INSERT INTO user_by_email (email, id, password_hash, created_at)
        VALUES (?, ?, ?, ?)`,
		"test@example.com", id, hashedPassword, createdAt,
	).Exec(); err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	jwtSecret := "testsecret"
	testHeader := login.LoginHandler(repo, jwtSecret)

	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testHeader(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", rr.Code)
	}

	var resp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Token == "" {
		t.Fatal("Expected JWT token in response, got empty string")
	}
}
