package config_test

import (
	"testing"

	"github.com/SuddenGunter/hsd/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadEnv_Success(t *testing.T) {
	// Set required environment variables
	t.Setenv("PORT", "8080")
	t.Setenv("MQTT_BROKER_HOST", "localhost")
	t.Setenv("MQTT_USERNAME", "testuser")
	t.Setenv("MQTT_PASSWORD", "testpass")
	t.Setenv("TELEGRAM_BOT_TOKEN", "123456:ABC-DEF1234")
	t.Setenv("TELEGRAM_CHAT_ID", "12345")

	cfg, err := config.LoadEnv()

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
	// Set required environment variables with custom MQTT port
	t.Setenv("PORT", "8080")
	t.Setenv("MQTT_BROKER_HOST", "mqtt.example.com")
	t.Setenv("MQTT_BROKER_PORT", "8883")
	t.Setenv("MQTT_USERNAME", "testuser")
	t.Setenv("MQTT_PASSWORD", "testpass")
	t.Setenv("TELEGRAM_BOT_TOKEN", "123456:ABC-DEF1234")
	t.Setenv("TELEGRAM_CHAT_ID", "12345")

	cfg, err := config.LoadEnv()

	require.NoError(t, err)
	assert.Equal(t, 8883, cfg.MQTT.BrokerPort)
}

func TestLoadEnv_WithZ2MDevices(t *testing.T) {
	// Set required environment variables with Z2M devices
	t.Setenv("PORT", "8080")
	t.Setenv("MQTT_BROKER_HOST", "localhost")
	t.Setenv("MQTT_USERNAME", "testuser")
	t.Setenv("MQTT_PASSWORD", "testpass")
	t.Setenv("Z2M_DEVICES", "device1,device2,device3")
	t.Setenv("TELEGRAM_BOT_TOKEN", "123456:ABC-DEF1234")
	t.Setenv("TELEGRAM_CHAT_ID", "12345")

	cfg, err := config.LoadEnv()

	require.NoError(t, err)
	assert.Equal(t, []string{"device1", "device2", "device3"}, cfg.Z2MDevices)
}

//nolint:paralleltest // Cannot use t.Parallel() with t.Setenv()
func TestLoadEnv_MissingRequiredFields(t *testing.T) {
	tests := []struct {
		name        string
		setupEnv    func(*testing.T)
		expectError bool
	}{
		{
			name: "missing PORT",
			setupEnv: func(t *testing.T) {
				t.Helper()
				t.Setenv("MQTT_BROKER_HOST", "localhost")
				t.Setenv("MQTT_USERNAME", "testuser")
				t.Setenv("MQTT_PASSWORD", "testpass")
				t.Setenv("TELEGRAM_BOT_TOKEN", "123456:ABC-DEF1234")
				t.Setenv("TELEGRAM_CHAT_ID", "12345")
			},
			expectError: true,
		},
		{
			name: "missing MQTT_BROKER_HOST",
			setupEnv: func(t *testing.T) {
				t.Helper()
				t.Setenv("PORT", "8080")
				t.Setenv("MQTT_USERNAME", "testuser")
				t.Setenv("MQTT_PASSWORD", "testpass")
				t.Setenv("TELEGRAM_BOT_TOKEN", "123456:ABC-DEF1234")
				t.Setenv("TELEGRAM_CHAT_ID", "12345")
			},
			expectError: true,
		},
		{
			name: "missing MQTT_USERNAME",
			setupEnv: func(t *testing.T) {
				t.Helper()
				t.Setenv("PORT", "8080")
				t.Setenv("MQTT_BROKER_HOST", "localhost")
				t.Setenv("MQTT_PASSWORD", "testpass")
				t.Setenv("TELEGRAM_BOT_TOKEN", "123456:ABC-DEF1234")
				t.Setenv("TELEGRAM_CHAT_ID", "12345")
			},
			expectError: true,
		},
		{
			name: "missing MQTT_PASSWORD",
			setupEnv: func(t *testing.T) {
				t.Helper()
				t.Setenv("PORT", "8080")
				t.Setenv("MQTT_BROKER_HOST", "localhost")
				t.Setenv("MQTT_USERNAME", "testuser")
				t.Setenv("TELEGRAM_BOT_TOKEN", "123456:ABC-DEF1234")
				t.Setenv("TELEGRAM_CHAT_ID", "12345")
			},
			expectError: true,
		},
		{
			name: "missing TELEGRAM_BOT_TOKEN",
			setupEnv: func(t *testing.T) {
				t.Helper()
				t.Setenv("PORT", "8080")
				t.Setenv("MQTT_BROKER_HOST", "localhost")
				t.Setenv("MQTT_USERNAME", "testuser")
				t.Setenv("MQTT_PASSWORD", "testpass")
				t.Setenv("TELEGRAM_CHAT_ID", "12345")
			},
			expectError: true,
		},
		{
			name: "missing TELEGRAM_CHAT_ID",
			setupEnv: func(t *testing.T) {
				t.Helper()
				t.Setenv("PORT", "8080")
				t.Setenv("MQTT_BROKER_HOST", "localhost")
				t.Setenv("MQTT_USERNAME", "testuser")
				t.Setenv("MQTT_PASSWORD", "testpass")
				t.Setenv("TELEGRAM_BOT_TOKEN", "123456:ABC-DEF1234")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		//nolint:paralleltest // Cannot use t.Parallel() with t.Setenv()
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv(t)

			cfg, err := config.LoadEnv()

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
	// Set environment variables with invalid PORT value
	t.Setenv("PORT", "invalid")
	t.Setenv("MQTT_BROKER_HOST", "localhost")
	t.Setenv("MQTT_USERNAME", "testuser")
	t.Setenv("MQTT_PASSWORD", "testpass")
	t.Setenv("TELEGRAM_BOT_TOKEN", "123456:ABC-DEF1234")
	t.Setenv("TELEGRAM_CHAT_ID", "12345")

	cfg, err := config.LoadEnv()

	require.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoadEnv_InvalidChatIDValue(t *testing.T) {
	// Set environment variables with invalid TELEGRAM_CHAT_ID value
	t.Setenv("PORT", "8080")
	t.Setenv("MQTT_BROKER_HOST", "localhost")
	t.Setenv("MQTT_USERNAME", "testuser")
	t.Setenv("MQTT_PASSWORD", "testpass")
	t.Setenv("TELEGRAM_BOT_TOKEN", "123456:ABC-DEF1234")
	t.Setenv("TELEGRAM_CHAT_ID", "invalid")

	cfg, err := config.LoadEnv()

	require.Error(t, err)
	assert.Nil(t, cfg)
}
