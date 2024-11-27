package handlers

import (
	"bytes"
	"encoding/json"
	"gomail/internal/config"
	"gomail/internal/models"
	"gomail/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostGameMail_Success(t *testing.T) {
	mockService := services.NewMailService(true)
	handler := NewMailHandler(mockService)
	config.LoadEnv()

	mailBody := models.MailBody{
		To:                 "recipient@example.com",
		RecipientVariables: map[string]map[string]string{"recipient@example.com": {"value": "Recipient"}},
		Game:               "Mafia",
	}

	body, _ := json.Marshal(mailBody)
	req, err := http.NewRequest("POST", "/mail/games", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.PostGameMail(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"message":"Mail sent successfully"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
