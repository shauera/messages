package main

import (
	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
)

func InitConfig() {
	// Source
	config.SetConfigName("config")
	config.AddConfigPath("/etc/messages")
	config.AddConfigPath(".")

	// Find and read the config file
	err := config.ReadInConfig()
	if err != nil {
		log.WithError(err).Fatal("Could not read config file")
	}

	// Defaults
	config.SetDefault(
		"logging", map[string]interface{}{
			"level": "info",
		},
	)

	config.SetDefault(
		"service", map[string]interface{}{
			"port":                  "8090",
			"shutdownGraceDuration": 10,
		},
	)

	config.SetDefault(
		"database", map[string]interface{}{
			"type": "memory",
		},
	)
}
