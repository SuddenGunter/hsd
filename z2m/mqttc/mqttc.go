package mqttc

import (
	"fmt"
	"time"

	"github.com/SuddenGunter/hsd/app/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Connect to the MQTT broker.
func Connect(cfg *config.Config) (mqtt.Client, error) {
	broker := cfg.MQTTBrokerIP
	port := cfg.MQTTBrokerPort
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("hsd")
	opts.SetUsername(cfg.MQTTUsername)
	opts.SetPassword(cfg.MQTTPassword)
	opts.SetConnectTimeout(10 * time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("mqtt: conn: %w", token.Error())
	}

	return client, nil
}
