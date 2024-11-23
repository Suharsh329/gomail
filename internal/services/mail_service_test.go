package services

import (
	"gomail/internal/config"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewMailService(t *testing.T) {
	service := NewMailService(true)
	if service == nil {
		t.Error("Expected MailService instance, got nil")
	}
}

func TestSendMail(t *testing.T) {
	service := NewMailService(true)

	config.LoadEnv()

	// Mock HTTP server
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Queued. Thank you."}`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// Call SendMail
	resp, err := service.SendMail("from", "to@example.com", "Test Subject", "Test Body", "mafia-game", map[string]map[string]string{"suharsh329": {"name": "Suharsh", "role": "Developer"}})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		return
	}

	// Check request body
	body := new(strings.Builder)
	_, err = io.Copy(body, resp.Body)
	if err != nil {
		t.Errorf("Expected no error reading response body, got %v", err)
	}
}
