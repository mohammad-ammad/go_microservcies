package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mohammad-ammad/auth-service/dto"
)

func TestLoginSuccess(t *testing.T) {
	defer CleanupTestDB()
	router := setupTestRouter()

	email := "test@example.com"
	password := "securepass"

	createTestUser(email, password)

	reqBody := dto.LoginRequest{
		Email:    email,
		Password: password,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if _, exists := resp["token"]; !exists {
		t.Errorf("Expected token in response, got: %v", resp)
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	defer CleanupTestDB()
	router := setupTestRouter()

	email := "test@example.com"
	password := "securepass"
	createTestUser(email, password)

	invalidLogin := dto.LoginRequest{
		Email:    email,
		Password: "wrongpass",
	}

	jsonBody, _ := json.Marshal(invalidLogin)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for invalid credentials, got %d", w.Code)
	}
}

func TestLoginNonExistentUser(t *testing.T) {
	defer CleanupTestDB()
	router := setupTestRouter()

	login := dto.LoginRequest{
		Email:    "notfound@example.com",
		Password: "irrelevant",
	}

	jsonBody, _ := json.Marshal(login)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for non-existent user, got %d", w.Code)
	}
}
