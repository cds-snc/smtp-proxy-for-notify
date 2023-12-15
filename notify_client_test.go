package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendEmail(t *testing.T) {
	// Create a mock NotifyClient
	client := &NotifyClient{
		Hostname: "http://example.com",
		ApiKey:   "test-api-key",
		Client:   &http.Client{},
	}

	// Create a mock NotifyEmail
	email := &NotifyEmail{
		TemplateId: "test-template-id",
		Personalisation: Body{
			Subject: "Test Subject",
			Body:    "Test Body",
		},
		Attachments: []Attachment{
			{
				File:          "test-file",
				Filename:      "test-filename",
				SendingMethod: "test-sending-method",
			},
		},
		Emails: []string{"test@example.com"},
	}

	// Create a mock response
	mockResponse := `{"status": "success"}`

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/notifications/email", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "ApiKey-v1 test-api-key", r.Header.Get("Authorization"))

		// Read the request body
		body, err := io.ReadAll(r.Body)
		assert.Nil(t, err)

		// Unmarshal the request body
		var requestPayload map[string]interface{}
		err = json.Unmarshal(body, &requestPayload)
		assert.Nil(t, err)

		// Verify the request payload
		assert.Equal(t, "test-template-id", requestPayload["template_id"])
		assert.Equal(t, "Test Subject", requestPayload["personalisation"].(map[string]interface{})["subject"])
		assert.Equal(t, "Test Body", requestPayload["personalisation"].(map[string]interface{})["body"])
		assert.Equal(t, "test-file", requestPayload["personalisation"].(map[string]interface{})["attachment_0"].(map[string]interface{})["file"])
		assert.Equal(t, "test-filename", requestPayload["personalisation"].(map[string]interface{})["attachment_0"].(map[string]interface{})["filename"])
		assert.Equal(t, "test-sending-method", requestPayload["personalisation"].(map[string]interface{})["attachment_0"].(map[string]interface{})["sending_method"])
		assert.Equal(t, "test@example.com", requestPayload["email_address"])

		// Write the mock response
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(mockResponse))
		assert.Nil(t, err)
	}))
	defer mockServer.Close()

	// Set the mock server URL as the NotifyClient hostname
	client.Hostname = mockServer.URL

	// Call the sendEmail function
	err := sendEmail(client, email)

	// Verify the result
	assert.Nil(t, err)
}
