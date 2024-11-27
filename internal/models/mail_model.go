package models

type MailBody struct {
	To                 string                       `json:"to"`
	RecipientVariables map[string]map[string]string `json:"recipient-variables"`
	Game               string                       `json:"game"`
}
