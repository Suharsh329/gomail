package models

type MailBody struct {
	From               string                       `json:"from"`
	To                 string                       `json:"to"`
	Subject            string                       `json:"subject"`
	Text               string                       `json:"text"`
	Template           string                       `json:"template"`
	RecipientVariables map[string]map[string]string `json:"recipient-variables"`
}
