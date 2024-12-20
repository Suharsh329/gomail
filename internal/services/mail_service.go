package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gomail/internal/config"
	"gomail/internal/models"
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

func getAPIKey() string {
	return config.GetEnvWithKey("MAILGUN_API_KEY", "")
}

func getDomainName() string {
	return config.GetEnvWithKey("MAILGUN_DOMAIN_NAME", "")
}

func getRoute() string {
	return config.GetEnvWithKey("MAILGUN_ROUTE", "")
}

func (s *MailService) CreateEmailBody(mailBody models.MailBody) (*multipart.Writer, bytes.Buffer, error) {
	// Create the POST body with multipart form data
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Add form fields
	_ = writer.WriteField("from", mailBody.From)
	_ = writer.WriteField("to", mailBody.To)
	_ = writer.WriteField("subject", mailBody.Subject)
	_ = writer.WriteField("text", mailBody.Text)
	_ = writer.WriteField("template", mailBody.Template)

	if mailBody.RecipientVariables != nil {
		// Convert recipient variables to JSON string
		recipientVariablesJSON, err := json.Marshal(mailBody.RecipientVariables)
		if err != nil {
			fmt.Println("Error marshalling recipient variables:", err)
			return nil, bytes.Buffer{}, err
		}
		_ = writer.WriteField("recipient-variables", string(recipientVariablesJSON))
	}

	if mailBody.Variables != nil {
		// Convert variables to JSON string
		variablesJSON, err := json.Marshal(mailBody.Variables)
		if err != nil {
			fmt.Println("Error marshalling variables:", err)
			return nil, bytes.Buffer{}, err
		}
		_ = writer.WriteField("t:variables", string(variablesJSON))
	}

	if s.IsTest {
		_ = writer.WriteField("o:testmode", "yes")
	}

	// Close the writer
	writer.Close()

	return writer, b, nil
}

func (s *MailService) CreateRequestBody(b bytes.Buffer, writer *multipart.Writer) (*http.Request, error) {
	mailgunApiKey := getAPIKey()
	mailgunRoute := getRoute()

	// Prepare the request
	req, err := http.NewRequest("POST", mailgunRoute, &b)
	if err != nil {
		return nil, err
	}

	// Add the necessary headers (including Basic Authentication)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth("api", mailgunApiKey)

	return req, nil
}

func (s *MailService) SendRequest(req *http.Request) (*http.Response, error) {
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

func (s *MailService) SendMail(mailBody models.MailBody) (*http.Response, error) {

	mailgunDomainName := getDomainName()

	if mailgunDomainName == "" {
		return nil, fmt.Errorf("mailgun domain name is empty")
	}

	if mailBody.Game != "" {
		mailBody.From = mailBody.Game + " - " + config.GetEnvWithKey("DOMAIN_NAME", "") + " <postmaster@" + mailgunDomainName + ">"

		switch mailBody.Game {
		case "Mafia":
			mailBody.Template = "mafia-game"
		case "Impostor":
			mailBody.Template = "impostor-game"
		default:
			mailBody.Template = "mafia-game"
		}
	} else {
		mailBody.From = config.GetEnvWithKey("DOMAIN_NAME", "") + " <postmaster@" + mailgunDomainName + ">"
	}

	writer, b, err := s.CreateEmailBody(mailBody)
	if err != nil {
		return nil, err
	}

	req, err := s.CreateRequestBody(b, writer)
	if err != nil {
		return nil, err
	}

	return s.SendRequest(req)
}
