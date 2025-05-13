package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mohammad-ammad/auth-service/dto"
)

func TestMeEndpoint(t *testing.T) {
	defer CleanupTestDB()
	router := setupTestRouter()

	email := "test@example.com"
	password := "securepass"

	createTestUser(email, password)

	loginReq := dto.LoginRequest{
		Email:    email,
		Password: password,
	}
	loginBody, _ := json.Marshal(loginReq)
	loginReqHTTP, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginBody))
	loginReqHTTP.Header.Set("Content-Type", "application/json")

	loginResp := httptest.NewRecorder()
	router.ServeHTTP(loginResp, loginReqHTTP)

	if loginResp.Code != http.StatusOK {
		t.Fatalf("Login failed, expected 200, got %d", loginResp.Code)
	}

	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)

	token, ok := loginData["token"].(string)
	if !ok || token == "" {
		t.Fatal("Token not returned in login response")
	}

	req, _ := http.NewRequest("GET", "/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 from /auth/me, got %d", w.Code)
	}

	var meResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &meResp); err != nil {
		t.Fatalf("Failed to parse /auth/me response: %v", err)
	}

	if _, ok := meResp["data"]; !ok {
		t.Errorf("Expected 'user' in response, got: %v", meResp)
	}
}
