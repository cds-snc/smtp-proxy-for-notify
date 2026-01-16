package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type Attachment struct {
	File          string `json:"file,omitempty"`
	Filename      string `json:"filename,omitempty"`
	SendingMethod string `json:"sending_method,omitempty"`
}

type Body struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type NotifyClient struct {
	ApiKey   string
	Client   *http.Client
	Hostname string
}

type NotifyEmail struct {
	Attachments     []Attachment `json:"-"`
	Emails          []string     `json:"-"`
	EmailAddress    string       `json:"email_address"`
	Personalisation Body         `json:"personalisation"`
	TemplateId      string       `json:"template_id"`
}

func newNotifyClient(apiKey string, hostname string) *NotifyClient {
	return &NotifyClient{
		ApiKey:   apiKey,
		Client:   &http.Client{Timeout: 10 * time.Second},
		Hostname: hostname,
	}
}

func sendEmail(client *NotifyClient, email *NotifyEmail) error {

	method := "POST"
	contentType := "application/json"
	resource := fmt.Sprintf("%s/v2/notifications/email", strings.Trim(client.Hostname, "/"))

	// Convert the struct to a map so we can join the attachments to the personalisation
	emailPayload := make(map[string]interface{})
	emailPayload["template_id"] = email.TemplateId
	emailPayload["personalisation"] = make(map[string]interface{})

	personalisation, ok := emailPayload["personalisation"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to convert personalisation to map[string]interface{}")
	}

	personalisation["subject"] = email.Personalisation.Subject
	personalisation["body"] = email.Personalisation.Body

	for i, attachment := range email.Attachments {
		key := fmt.Sprintf("attachment_%d", i)
		personalisation[key] = map[string]interface{}{
			"file":           attachment.File,
			"filename":       attachment.Filename,
			"sending_method": attachment.SendingMethod,
		}
	}

	for _, email_address := range email.Emails {

		emailPayload["email_address"] = email_address

		body, err := json.Marshal(emailPayload)

		log.Info().Msgf("Sending email to : %s", email_address)
		if err != nil {
			return err
		}

		req, err := http.NewRequest(method, resource, bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", fmt.Sprintf("ApiKey-v1 %s", client.ApiKey))

		if err := doSendEmail(client.Client, req); err != nil {
			return err
		}
	}

	return nil
}

func doSendEmail(client *http.Client, req *http.Request) error {
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Msgf("Error sending email: %s", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		return nil
	}

	// A non-201 status code so craft the error.
	log.Error().Msgf("Unexpected status code: %d", resp.StatusCode)
	respbody, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Error().Msgf("Error reading response body: %s", err)
		return err
	}
	log.Error().Msgf("Response: %s", respbody)

	return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
