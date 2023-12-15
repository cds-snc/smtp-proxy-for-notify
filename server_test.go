package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthPlain_Passes(t *testing.T) {
	// Create a mock Session
	session := Session{
		Authenticated: false,
		Config: &Config{
			Smtp: struct {
				Hostname    string
				Port        int
				Username    string
				Password    string
				UseTLS      bool
				TlsCertFile string
				TlsKeyFile  string
			}{
				Username: "test-username",
				Password: "test-password",
			},
		},
	}

	// Call the AuthPlain method
	err := session.AuthPlain("test-username", "test-password")

	// Verify the result
	assert.Equal(t, session.Authenticated, true)
	assert.Nil(t, err)
}

func TestAuthPlain_Fails(t *testing.T) {
	// Create a mock Session
	session := Session{
		Authenticated: false,
		Config: &Config{
			Smtp: struct {
				Hostname    string
				Port        int
				Username    string
				Password    string
				UseTLS      bool
				TlsCertFile string
				TlsKeyFile  string
			}{
				Username: "test-username",
				Password: "test-password",
			},
		},
	}

	// Call the AuthPlain method
	err := session.AuthPlain("test-username", "test-password-1")

	// Verify the result
	assert.Equal(t, session.Authenticated, false)
	assert.Equal(t, err.Error(), "invalid username or password")
}

func TestSession_Data(t *testing.T) {
	// Create a mock Session
	session := Session{
		Authenticated: true,
		Email: &NotifyEmail{
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
			Emails: []string{"test@test.com"},
		},
		Config: &Config{
			Notify: struct {
				ApiKey     string
				Hostname   string
				TemplateId string
			}{
				ApiKey:   "test-api-key",
				Hostname: "http://example.com",
			},
		},
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
		assert.Equal(t, "Test Body\r", requestPayload["personalisation"].(map[string]interface{})["body"])
		assert.Equal(t, "test-file", requestPayload["personalisation"].(map[string]interface{})["attachment_0"].(map[string]interface{})["file"])
		assert.Equal(t, "test-filename", requestPayload["personalisation"].(map[string]interface{})["attachment_0"].(map[string]interface{})["filename"])
		assert.Equal(t, "test-sending-method", requestPayload["personalisation"].(map[string]interface{})["attachment_0"].(map[string]interface{})["sending_method"])
		assert.Contains(t, [3]string{"test@test.com", "test1@test.com", "test2@test.com"}, requestPayload["email_address"])

		// Write the mock response
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(mockResponse))
		assert.Nil(t, err)
	}))
	defer mockServer.Close()

	// Set the mock server URL as the NotifyClient hostname
	session.Config.Notify.Hostname = mockServer.URL

	email := "" +
		"To: <test@test.com>\r\n" +
		"Cc: <test1@test.com>\r\n" +
		"Bcc: <test2@test.com>\r\n" +
		"Subject: Test Subject\r\n" +
		"\r\n" +
		"Test Body\r\n"

	// Call the Data method
	err := session.Data(strings.NewReader(email))

	// Verify the result
	assert.Nil(t, err)
}

func TestSession_Logout(t *testing.T) {
	// Create a mock Session
	session := Session{}

	// Call the Logout method
	err := session.Logout()

	// Verify the result
	assert.Equal(t, session.Authenticated, false)
	assert.Nil(t, err)
}

func TestSession_MailWithoutAuth(t *testing.T) {
	// Create a mock Session
	session := Session{
		Authenticated: false,
	}

	// Call the Mail method
	err := session.Mail("test@test.com", nil)

	// Verify the result
	assert.Equal(t, session.Authenticated, false)
	assert.Equal(t, err.Error(), "not authenticated")
}

func TestSession_MailWithAuth(t *testing.T) {
	// Create a mock Session
	session := Session{
		Authenticated: true,
	}

	// Call the Mail method
	err := session.Mail("test@test.com", nil)

	// Verify the result
	assert.Equal(t, session.Authenticated, true)
	assert.Nil(t, err)
}

func TestSession_RcptWithoutAuth(t *testing.T) {
	// Create a mock Session
	session := Session{
		Authenticated: false,
	}

	// Call the Rcpt method
	err := session.Rcpt("test@test.com", nil)

	// Verify the result
	assert.Equal(t, session.Authenticated, false)
	assert.Equal(t, err.Error(), "not authenticated")
}

func TestSession_RcptWithAuth(t *testing.T) {
	// Create a mock Session
	session := Session{
		Authenticated: true,
		Email: &NotifyEmail{
			Emails: []string{},
		},
	}

	// Call the Rcpt method
	err := session.Rcpt("test@test.com", nil)

	// Verify the result
	assert.Equal(t, session.Authenticated, true)
	assert.Equal(t, session.Email.Emails, []string{"test@test.com"})
	assert.Nil(t, err)
}

func TestSession_Reset(t *testing.T) {
	// Create a mock Session
	session := Session{
		Authenticated: true,
		Config: &Config{
			Notify: struct {
				ApiKey     string
				Hostname   string
				TemplateId string
			}{
				TemplateId: "template-id",
			}},
		Email: &NotifyEmail{
			Emails: []string{"test@test.com"},
		},
	}

	// Call the Reset method
	session.Reset()

	// Verify the result
	assert.Equal(t, session.Authenticated, false)
	assert.Equal(t, session.Email.TemplateId, "template-id")
	assert.Equal(t, session.Email.Emails, []string(nil))
}

func Test_startSmtpServer_withoutTls(t *testing.T) {
	// Create a mock Backend
	config := Config{
		Smtp: struct {
			Hostname    string
			Port        int
			Username    string
			Password    string
			UseTLS      bool
			TlsCertFile string
			TlsKeyFile  string
		}{
			Hostname: "localhost",
			Port:     2525,
			Username: "test-username",
			Password: "test-password",
			UseTLS:   false,
		},
	}

	// Call the startSmtpServer method
	go func() {
		startSmtpServer(&config)
	}()
}

func Test_startSmtpServer_withTls(t *testing.T) {
	// Create a mock Backend
	config := Config{
		Smtp: struct {
			Hostname    string
			Port        int
			Username    string
			Password    string
			UseTLS      bool
			TlsCertFile string
			TlsKeyFile  string
		}{
			Hostname:    "localhost",
			Port:        2525,
			Username:    "test-username",
			Password:    "test-password",
			UseTLS:      true,
			TlsCertFile: "./example_certs/server.crt",
			TlsKeyFile:  "./example_certs/server.key",
		},
	}

	// Call the startSmtpServer method
	go func() {
		startSmtpServer(&config)
	}()
}
