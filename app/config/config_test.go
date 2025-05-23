package config_test

import (
	"testing"

	"github.com/SuddenGunter/hsd/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadEnv_Success(t *testing.T) {
	t.Parallel()

	envMap := map[string]string{
		"PORT":               "8080",
		"MQTT_BROKER_HOST":   "localhost",
		"MQTT_USERNAME":      "testuser",
		"MQTT_PASSWORD":      "testpass",
		"TELEGRAM_BOT_TOKEN": "123456:ABC-DEF1234",
		"TELEGRAM_CHAT_ID":   "12345",
	}

	cfg, err := config.LoadFromEnvMap(envMap)

	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "localhost", cfg.MQTT.BrokerHost)
	assert.Equal(t, 1883, cfg.MQTT.BrokerPort) // default value
	assert.Equal(t, "testuser", cfg.MQTT.Username)
	assert.Equal(t, "testpass", cfg.MQTT.Password)
	assert.Equal(t, "123456:ABC-DEF1234", cfg.Telegram.BotToken)
	assert.Equal(t, int64(12345), cfg.Telegram.ChatID)
	assert.Empty(t, cfg.Z2MDevices) // optional field
}

func TestLoadEnv_WithCustomMQTTPort(t *testing.T) {
	t.Parallel()

	envMap := map[string]string{
		"PORT":               "8080",
		"MQTT_BROKER_HOST":   "mqtt.example.com",
		"MQTT_BROKER_PORT":   "8883",
		"MQTT_USERNAME":      "testuser",
		"MQTT_PASSWORD":      "testpass",
		"TELEGRAM_BOT_TOKEN": "123456:ABC-DEF1234",
		"TELEGRAM_CHAT_ID":   "12345",
	}

	cfg, err := config.LoadFromEnvMap(envMap)

	require.NoError(t, err)
	assert.Equal(t, 8883, cfg.MQTT.BrokerPort)
}

func TestLoadEnv_WithZ2MDevices(t *testing.T) {
	t.Parallel()

	envMap := map[string]string{
		"PORT":               "8080",
		"MQTT_BROKER_HOST":   "localhost",
		"MQTT_USERNAME":      "testuser",
		"MQTT_PASSWORD":      "testpass",
		"Z2M_DEVICES":        "device1,device2,device3",
		"TELEGRAM_BOT_TOKEN": "123456:ABC-DEF1234",
		"TELEGRAM_CHAT_ID":   "12345",
	}

	cfg, err := config.LoadFromEnvMap(envMap)

	require.NoError(t, err)
	assert.Equal(t, []string{"device1", "device2", "device3"}, cfg.Z2MDevices)
}

func TestLoadEnv_MissingRequiredFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupEnv    func() map[string]string
		expectError bool
	}{
		{
			name: "missing PORT",
			setupEnv: func() map[string]string {
				return map[string]string{
					"MQTT_BROKER_HOST":   "localhost",
					"MQTT_USERNAME":      "testuser",
					"MQTT_PASSWORD":      "testpass",
					"TELEGRAM_BOT_TOKEN": "123456:ABC-DEF1234",
					"TELEGRAM_CHAT_ID":   "12345",
				}
			},
			expectError: true,
		},
		{
			name: "missing MQTT_BROKER_HOST",
			setupEnv: func() map[string]string {
				return map[string]string{
					"PORT":               "8080",
					"MQTT_USERNAME":      "testuser",
					"MQTT_PASSWORD":      "testpass",
					"TELEGRAM_BOT_TOKEN": "123456:ABC-DEF1234",
					"TELEGRAM_CHAT_ID":   "12345",
				}
			},
			expectError: true,
		},
		{
			name: "missing MQTT_USERNAME",
			setupEnv: func() map[string]string {
				return map[string]string{
					"PORT":               "8080",
					"MQTT_BROKER_HOST":   "localhost",
					"MQTT_PASSWORD":      "testpass",
					"TELEGRAM_BOT_TOKEN": "123456:ABC-DEF1234",
					"TELEGRAM_CHAT_ID":   "12345",
				}
			},
			expectError: true,
		},
		{
			name: "missing MQTT_PASSWORD",
			setupEnv: func() map[string]string {
				return map[string]string{
					"PORT":               "8080",
					"MQTT_BROKER_HOST":   "localhost",
					"MQTT_USERNAME":      "testuser",
					"TELEGRAM_BOT_TOKEN": "123456:ABC-DEF1234",
					"TELEGRAM_CHAT_ID":   "12345",
				}
			},
			expectError: true,
		},
		{
			name: "missing TELEGRAM_BOT_TOKEN",
			setupEnv: func() map[string]string {
				return map[string]string{
					"PORT":             "8080",
					"MQTT_BROKER_HOST": "localhost",
					"MQTT_USERNAME":    "testuser",
					"MQTT_PASSWORD":    "testpass",
					"TELEGRAM_CHAT_ID": "12345",
				}
			},
			expectError: true,
		},
		{
			name: "missing TELEGRAM_CHAT_ID",
			setupEnv: func() map[string]string {
				return map[string]string{
					"PORT":               "8080",
					"MQTT_BROKER_HOST":   "localhost",
					"MQTT_USERNAME":      "testuser",
					"MQTT_PASSWORD":      "testpass",
					"TELEGRAM_BOT_TOKEN": "123456:ABC-DEF1234",
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			envMap := tt.setupEnv()
			cfg, err := config.LoadFromEnvMap(envMap)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, cfg)
			}
		})
	}
}

func TestLoadEnv_InvalidPortValue(t *testing.T) {
	t.Parallel()

	envMap := map[string]string{
		"PORT":               "invalid",
		"MQTT_BROKER_HOST":   "localhost",
		"MQTT_USERNAME":      "testuser",
		"MQTT_PASSWORD":      "testpass",
		"TELEGRAM_BOT_TOKEN": "123456:ABC-DEF1234",
		"TELEGRAM_CHAT_ID":   "12345",
	}

	cfg, err := config.LoadFromEnvMap(envMap)

	require.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoadEnv_InvalidChatIDValue(t *testing.T) {
	t.Parallel()

	envMap := map[string]string{
		"PORT":               "8080",
		"MQTT_BROKER_HOST":   "localhost",
		"MQTT_USERNAME":      "testuser",
		"MQTT_PASSWORD":      "testpass",
		"TELEGRAM_BOT_TOKEN": "123456:ABC-DEF1234",
		"TELEGRAM_CHAT_ID":   "invalid",
	}

	cfg, err := config.LoadFromEnvMap(envMap)

	require.Error(t, err)
	assert.Nil(t, cfg)
}
