package main

import "github.com/rs/zerolog/log"

func main() {
	// Initialize configuration
	config, err := initConfig()

	if err != nil {
		log.Fatal().Msgf("Error initializing configuration: %s", err)
	}

	// Start the SMTP server
	startSmtpServer(config)
}
