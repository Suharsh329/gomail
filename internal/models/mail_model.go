package models

type MailBody struct {
	From               string                       `json:"from"`
	Game               string                       `json:"game"`
	RecipientVariables map[string]map[string]string `json:"recipient-variables"`
	Subject            string                       `json:"subject"`
	Template           string                       `json:"template"`
	Text               string                       `json:"text"`
	To                 string                       `json:"to"`
	Variables          map[string]string            `json:"variables"`
}
