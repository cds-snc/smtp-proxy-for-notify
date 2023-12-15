package main

import (
	"crypto/tls"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/DusanKasan/parsemail"
	"github.com/emersion/go-smtp"
	"github.com/rs/zerolog/log"
)

type Backend struct {
	Config *Config
}

func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		Authenticated: false,
		Config:        bkd.Config,
		Email: &NotifyEmail{
			TemplateId: bkd.Config.Notify.TemplateId,
		},
	}, nil
}

type Session struct {
	Authenticated bool
	Config        *Config
	Email         *NotifyEmail
}

func (s *Session) AuthPlain(username, password string) error {
	if username != s.Config.Smtp.Username || password != s.Config.Smtp.Password {
		log.Error().Msgf("Invalid username or password: %s", username)
		s.Authenticated = false
		s.Logout()
		return errors.New("invalid username or password")
	}
	log.Info().Msgf("User %s logged in", username)
	s.Authenticated = true
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	if !s.Authenticated {
		s.Logout()
		return errors.New("not authenticated")
	}
	log.Info().Msgf("Mail from: %s", from)
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if !s.Authenticated {
		s.Logout()
		return errors.New("not authenticated")
	}
	log.Info().Msgf("Rcpt to: %s", to)
	s.Email.Emails = append(s.Email.Emails, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if !s.Authenticated {
		return errors.New("not authenticated")
	}

	if data, err := io.ReadAll(r); err != nil {
		return err
	} else {
		// Parse the email
		email, err := parsemail.Parse(strings.NewReader(string(data)))

		if err != nil {
			log.Error().Msgf("Error parsing email: %s", err)
			return err
		}

		// Add cc emails
		for _, address := range email.Cc {
			s.Email.Emails = append(s.Email.Emails, address.Address)
		}

		// Add bcc emails
		for _, address := range email.Bcc {
			s.Email.Emails = append(s.Email.Emails, address.Address)
		}

		s.Email.Personalisation.Subject = email.Subject
		s.Email.Personalisation.Body = email.TextBody

		// Add attachments
		for _, attachment := range email.Attachments {
			attachment_data, err := io.ReadAll(attachment.Data)
			if err != nil {
				return err
			}
			s.Email.Attachments = append(s.Email.Attachments, Attachment{
				File:          b64.StdEncoding.EncodeToString(attachment_data),
				Filename:      attachment.Filename,
				SendingMethod: "attach",
			})
		}

		client := newNotifyClient(s.Config.Notify.ApiKey, s.Config.Notify.Hostname)
		if err := sendEmail(client, s.Email); err != nil {
			return err
		}
	}
	return nil
}

func (s *Session) Reset() {
	s.Authenticated = false
	s.Email = new(NotifyEmail)
	s.Email.TemplateId = s.Config.Notify.TemplateId
}

func (s *Session) Logout() error {
	s.Authenticated = false
	return nil
}

func startSmtpServer(config *Config) {
	backend := &Backend{
		Config: config,
	}

	s := smtp.NewServer(backend)

	s.Addr = fmt.Sprintf("%s:%d", config.Smtp.Hostname, config.Smtp.Port)
	s.Domain = config.Smtp.Hostname
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 10485760
	s.MaxRecipients = 10

	if config.Smtp.UseTLS {
		s.AllowInsecureAuth = false
		s.EnableREQUIRETLS = true

		cer, err := tls.LoadX509KeyPair(config.Smtp.TlsCertFile, config.Smtp.TlsKeyFile)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to load TLS certificate")
		}
		s.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cer}}

		log.Info().Msgf("SMTP server listening with TLS at %s", s.Addr)
		if err := s.ListenAndServeTLS(); err != nil {
			log.Fatal().Err(err).Msg("SMTP server failed with TLS")
		}
	} else {
		s.AllowInsecureAuth = true
		log.Warn().Msg("SMTP server listening without TLS! DO NOT USE IN PRODUCTION!")
		log.Info().Msgf("SMTP server listening on %s", s.Addr)
		if err := s.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("SMTP server failed")
		}
	}
}
