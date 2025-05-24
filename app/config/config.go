package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config represents the configuration of the app.
type Config struct {
	Port int `env:"PORT,required"`

	MQTT mqttConfig `envPrefix:"MQTT_"`

	Z2MDevices []string `env:"Z2M_DEVICES"`

	Telegram telegramConfig `envPrefix:"TELEGRAM_"`
}

type mqttConfig struct {
	BrokerHost string `env:"BROKER_HOST,required"`
	BrokerPort int    `env:"BROKER_PORT" envDefault:"1883"`
	Username   string `env:"USERNAME,required"`
	Password   string `env:"PASSWORD,required"`
}

type telegramConfig struct {
	BotToken string `env:"BOT_TOKEN,required"`
	ChatID   int64  `env:"CHAT_ID,required"`
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
