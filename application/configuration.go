package application

import (
	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"

	"strings"
)

//InitConfig - initialize message manager instance defaults
//This will read the configuration file and set defaults for each
//missing configuration entry
func InitConfig() {
	// Prefix all envirinment variables with "messages"
	config.SetEnvPrefix("messages")
	// Compound variable names with `_` instead of `.``
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// Automatically read (and use as overrides) environment variables
	config.AutomaticEnv()

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
			"healthPort":            "8091",
			"healthScanInterval":    "30s",
			"shutdownGraceDuration": "10s",
		},
	)

	config.SetDefault(
		"database", map[string]interface{}{
			"type": "memory",
		},
	)
}
