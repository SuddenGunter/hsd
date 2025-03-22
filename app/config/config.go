package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config represents the configuration of the app.
type Config struct {
	Port int `env:"PORT,required"`

	MQTTBrokerHost string `env:"MQTT_BROKER_HOST,required"`
	MQTTBrokerPort int    `env:"MQTT_BROKER_PORT" envDefault:"1883"`
	MQTTUsername   string `env:"MQTT_USERNAME,required"`
	MQTTPassword   string `env:"MQTT_PASSWORD,required"`

	Z2MDevices []string `env:"Z2M_DEVICES"`

	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
	TelegramChatID   int64  `env:"TELEGRAM_CHAT_ID,required"`
}

// LoadEnv loads the configuration from the environment.
func LoadEnv() (*Config, error) {
	cfg := Config{}

	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("env: %w", err)
	}

	return &cfg, nil
}
