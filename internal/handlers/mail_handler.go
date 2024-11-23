package handlers

import (
	"encoding/json"
	"gomail/internal/models"
	"gomail/internal/services"
	"gomail/internal/utils"
	"net/http"
)

type MailHandler struct {
	Service *services.MailService
}

func NewMailHandler(service *services.MailService) *MailHandler {
	return &MailHandler{Service: service}
}

func (h *MailHandler) PostMail(w http.ResponseWriter, r *http.Request) {
	var mailBody models.MailBody

	json.NewDecoder(r.Body).Decode(&mailBody)

	_, err := h.Service.SendMail(mailBody.From, mailBody.To, mailBody.Subject, mailBody.Text, mailBody.Template, mailBody.RecipientVariables)

	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Response(w, http.StatusOK, map[string]string{"message": "Mail sent successfully"})
}
