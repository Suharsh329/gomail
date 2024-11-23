package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gomail/internal/config"
	"log"
	"mime/multipart"
	"net/http"
)

type MailService struct {
	IsTest bool
}

func NewMailService(isTest bool) *MailService {
	return &MailService{IsTest: isTest}
}

func getEnvVariables() (string, string, string, error) {
	get := config.GetEnvWithKey
	mailgunApiKey := get("MAILGUN_API_KEY", "")
	mailgunDomainName := get("MAILGUN_DOMAIN_NAME", "")
	mailgunRoute := get("MAILGUN_ROUTE", "https://api.mailgun.net/v3/"+mailgunDomainName+"/messages")

	if mailgunApiKey == "" {
		return "", "", "", fmt.Errorf("mailgun api key is empty")
	}

	if mailgunDomainName == "" {
		return "", "", "", fmt.Errorf("mailgun domain name is empty")
	}

	return mailgunApiKey, mailgunDomainName, mailgunRoute, nil
}

func (s *MailService) SendMail(_from, _to, _subject, _text, template string, _recipientVariables map[string]map[string]string) (*http.Response, error) {

	mailgunApiKey, mailgunDomainName, mailgunRoute, err := getEnvVariables()

	if err != nil {
		return nil, err
	}

	from := _from + " <postmaster@" + mailgunDomainName + ">"
	to := _to
	subject := _subject
	text := _text

	// Define recipient variables
	recipientVariables := _recipientVariables

	// Convert recipient variables to JSON string
	recipientVariablesJSON, err := json.Marshal(recipientVariables)
	if err != nil {
		fmt.Println("Error marshalling recipient variables:", err)
		return nil, err
	}

	// Create the POST body with multipart form data
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Add form fields
	_ = writer.WriteField("from", from)
	_ = writer.WriteField("to", to)
	_ = writer.WriteField("subject", subject)
	_ = writer.WriteField("text", text)
	_ = writer.WriteField("template", template)
	_ = writer.WriteField("recipient-variables", string(recipientVariablesJSON))

	if s.IsTest {
		_ = writer.WriteField("o:testmode", "yes")
	}

	// Close the writer
	writer.Close()

	// Prepare the request
	req, err := http.NewRequest("POST", mailgunRoute, &b)
	if err != nil {
		return nil, err
	}

	// Add the necessary headers (including Basic Authentication)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth("api", mailgunApiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Mailgun response: %v", resp)
		return nil, fmt.Errorf("mailgun request failed with status: %v", resp.Status)
	}

	return resp, nil
}
