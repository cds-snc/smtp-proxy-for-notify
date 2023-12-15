package main

import (
	"errors"
	"regexp"

	"github.com/spf13/viper"
)

type Config struct {
	// Notify settings
	Notify struct {
		ApiKey     string
		Hostname   string
		TemplateId string
	}

	// SMTP settings
	Smtp struct {
		// Hostname to listen on
		Hostname string
		Port     int

		// Username and password for authentication
		Username string
		Password string

		// Start TLS
		UseTLS      bool
		TlsCertFile string
		TlsKeyFile  string
	}
}

func initConfig() (*Config, error) {
	viper.AutomaticEnv()

	var configuration Config
	var err error

	// Set default values
	viper.SetDefault("LogLevel", "info")
	viper.SetDefault("Notify_ApiKey", "")
	viper.SetDefault("Notify_Hostname", "https://api.notification.canada.ca")
	viper.SetDefault("Notify_Template_Id", "")
	viper.SetDefault("Smtp_Hostname", "localhost")
	viper.SetDefault("Smtp_Port", 1025)
	viper.SetDefault("Smtp_Username", "")
	viper.SetDefault("Smtp_Password", "")
	viper.SetDefault("Smtp_Use_tls", false)
	viper.SetDefault("Smtp_tls_cert_file", "")
	viper.SetDefault("Smtp_tls_key_file", "")

	configuration.Notify.ApiKey = viper.GetString("Notify_ApiKey")
	configuration.Notify.Hostname = viper.GetString("Notify_Hostname")
	configuration.Notify.TemplateId = viper.GetString("Notify_Template_Id")
	configuration.Smtp.Hostname = viper.GetString("Smtp_Hostname")
	configuration.Smtp.Port = viper.GetInt("Smtp_Port")
	configuration.Smtp.Username = viper.GetString("Smtp_Username")
	configuration.Smtp.Password = viper.GetString("Smtp_Password")
	configuration.Smtp.UseTLS = viper.GetBool("Smtp_Use_TLS")
	configuration.Smtp.TlsCertFile = viper.GetString("Smtp_tls_cert_file")
	configuration.Smtp.TlsKeyFile = viper.GetString("Smtp_tls_key_file")

	// Validate username is no less than three characters
	if len(configuration.Smtp.Username) < 3 {
		err := errors.New("username must be at least three characters")
		return &configuration, err
	}

	// Validate password is no less than fourteen characters
	if len(configuration.Smtp.Password) < 14 {
		err := errors.New("password must be at least fourteen characters")
		return &configuration, err
	}

	// Validate API key starts with gcntfy and is not less than 81 characters
	if len(configuration.Notify.ApiKey) < 81 || configuration.Notify.ApiKey[:6] != "gcntfy" {
		err := errors.New("API key must start with gcntfy and be at least 81 characters")
		return &configuration, err
	}

	// Validate Notify Template ID matches a UUIDv4 using regex
	r, _ := regexp.Compile(`^[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-4[0-9a-fA-F]{3}\-[89abAB][0-9a-fA-F]{3}\-[0-9a-fA-F]{12}$`)
	if !r.MatchString(configuration.Notify.TemplateId) {
		err := errors.New("notify Template ID must be a UUIDv4")
		return &configuration, err
	}

	// If TLS is enabled, validate the certificate and key file paths
	if configuration.Smtp.UseTLS {
		if configuration.Smtp.TlsCertFile == "" {
			err := errors.New("TLS certificate file path must be specified")
			return &configuration, err
		}
		if configuration.Smtp.TlsKeyFile == "" {
			err := errors.New("TLS key file path must be specified")
			return &configuration, err
		}
	}

	return &configuration, err
}
