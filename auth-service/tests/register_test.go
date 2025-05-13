package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mohammad-ammad/auth-service/dto"
)

func TestRegisterSuccess(t *testing.T) {
	defer CleanupTestDB()

	router := setupTestRouter()

	body := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "securepass",
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestRegisterValidationError(t *testing.T) {
	defer CleanupTestDB()

	router := setupTestRouter()

	body := dto.RegisterRequest{
		Username: "",
		Email:    "invalid-email",
		Password: "123",
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
