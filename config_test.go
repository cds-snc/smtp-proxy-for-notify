package main

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	viper.AutomaticEnv()

	// Test case 1: Valid configuration
	_, err := initConfig()
	assert.Nil(t, err)

	// Test case 2: Invalid username
	viper.Set("Smtp_Username", "ab")
	_, err = initConfig()
	assert.Equal(t, "username must be at least three characters", err.Error())

	// Test case 3: Invalid password
	viper.Set("Smtp_Username", "valid_username")
	viper.Set("Smtp_Password", "short")
	_, err = initConfig()
	assert.Equal(t, "password must be at least fourteen characters", err.Error())

	// Test case 4: Invalid API key
	viper.Set("Smtp_Password", "valid_password")
	viper.Set("Notify_ApiKey", "invalid_key")
	_, err = initConfig()
	assert.Equal(t, "API key must start with gcntfy and be at least 81 characters", err.Error())

	// Test case 5: Invalid Notify Template ID
	viper.Set("Notify_ApiKey", "gcntfy-test-00000000-0000-4000-8000-000000000000-00000000-0000-4000-8000-000000000000")
	viper.Set("Notify_Template_Id", "invalid_template_id")
	_, err = initConfig()
	assert.Equal(t, "notify Template ID must be a UUIDv4", err.Error())

	// Test case 6: Invalid TLS configuration for
	viper.Set("Notify_Template_Id", "00000000-0000-4000-8000-000000000000")
	viper.Set("Smtp_Use_tls", true)
	viper.Set("Smtp_tls_cert_file", "")
	viper.Set("Smtp_tls_key_file", "")
	_, err = initConfig()
	assert.Equal(t, "TLS certificate file path must be specified", err.Error())

	// Test case 7: Invalid TLS configuration for
	viper.Set("Smtp_tls_cert_file", "valid_cert_file")
	viper.Set("Smtp_tls_key_file", "")
	_, err = initConfig()
	assert.Equal(t, "TLS key file path must be specified", err.Error())
}
