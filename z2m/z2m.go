package z2m

import (
	"context"
	"log/slog"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Msg represents a message from zigbee2mqtt.
type Msg struct {
	Payload []byte
	Device  string
}

type msgHandler interface {
	Handle(ctx context.Context, msg Msg)
}

// Zigbee2MQTTListener listens to zigbee2mqtt messages and forwards them to respective handlers.
type Zigbee2MQTTListener struct {
	client              mqtt.Client
	dataHandler         msgHandler
	availabilityHandler msgHandler
	allowedDevices      map[string]struct{}

	l *slog.Logger
}

// NewZigbee2MQTTListener returns a new Zigbee2MQTTListener.
func NewZigbee2MQTTListener(
	client mqtt.Client,
	dataHandler msgHandler,
	availabilityHandler msgHandler,
	allowedDevices []string,
	l *slog.Logger,
) *Zigbee2MQTTListener {
	m := make(map[string]struct{}, len(allowedDevices))
	for _, d := range allowedDevices {
		m[d] = struct{}{}
	}

	if len(m) == 0 {
		l.Error("no devices were enabled for zigbee2mqtt listener")
	}

	return &Zigbee2MQTTListener{client: client, dataHandler: dataHandler, availabilityHandler: availabilityHandler, allowedDevices: m, l: l}
}

// Subscribe to the zigbee2mqtt/# topic.
func (listener *Zigbee2MQTTListener) Subscribe() {
	topic := "zigbee2mqtt/#"
	token := listener.client.Subscribe(topic, 1, listener.onMessage)
	token.Wait()
}

func (listener *Zigbee2MQTTListener) onMessage(_ mqtt.Client, msg mqtt.Message) {
	defer msg.Ack()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	topic := msg.Topic()
	topic = strings.TrimPrefix(topic, "zigbee2mqtt/")

	if strings.HasPrefix(topic, "bridge") {
		listener.l.Debug("bridge state msg received, skip")

		return
	}

	if strings.HasSuffix(topic, "/availability") {
		device := strings.TrimSuffix(topic, "/availability")
		if _, ok := listener.allowedDevices[device]; !ok {
			listener.l.Debug("device not allowed", "device", device)

			return
		}

		listener.availabilityHandler.Handle(ctx, Msg{
			Payload: msg.Payload(),
			Device:  device,
		})

		return
	}

	if _, ok := listener.allowedDevices[topic]; !ok {
		listener.l.Debug("device not allowed", "device", topic)

		return
	}

	listener.dataHandler.Handle(ctx, Msg{
		Payload: msg.Payload(),
		Device:  topic,
	})
}
