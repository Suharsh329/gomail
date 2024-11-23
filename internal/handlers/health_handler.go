package handlers

import (
	"gomail/internal/utils"
	"net/http"
)

func NewHealthHandler() *MailHandler {
	return &MailHandler{}
}

func (h *MailHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.Response(w, http.StatusOK, map[string]string{"status": "ok"})
}
