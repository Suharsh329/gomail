package handlers

import (
	"encoding/json"
	"gomail/internal/models"
	"gomail/internal/services"
	"gomail/internal/utils"
	"log"
	"net/http"
)

type MailHandler struct {
	Service *services.MailService
}

func NewMailHandler(service *services.MailService) *MailHandler {
	return &MailHandler{Service: service}
}

func (h *MailHandler) PostGameMail(w http.ResponseWriter, r *http.Request) {
	var mailBody models.MailBody

	err := json.NewDecoder(r.Body).Decode(&mailBody)
	if err != nil {
		log.Printf("Error decoding mail body: %v", err)
		utils.Response(w, http.StatusBadRequest, "Bad request")
		return
	}

	from := mailBody.Game + "- Commando Lizard"

	var template string

	switch mailBody.Game {
	case "Mafia":
		template = "mafia-game"
	case "Impostor":
		template = "impostor-game"
	default:
		template = "mafia-game"
	}

	if mailBody.Game == "Impostor" {
		_, err = h.Service.SendMailWithVariables(from, mailBody.To, "", "", template, mailBody.RecipientVariables)
	} else {
		_, err = h.Service.SendMail(from, mailBody.To, "", "", template, mailBody.RecipientVariables)
	}
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	utils.Response(w, http.StatusOK, map[string]string{"message": "Mail sent successfully"})
}
