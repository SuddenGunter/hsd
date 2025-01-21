package device

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/SuddenGunter/hsd/z2m"
)

// DataHandler handles payload messages from zigbee2mqtt about device state updates.
// Currently it only handles door sensor messages.
type DataHandler struct {
	deviceNotifier deviceNotifier
	l              *slog.Logger
}

// NewDataHandler returns a new DataHandler.
func NewDataHandler(deviceNotifier deviceNotifier, l *slog.Logger) *DataHandler {
	return &DataHandler{
		deviceNotifier: deviceNotifier,
		l:              l,
	}
}

// Handle payload messages from zigbee2mqtt about device state updates.
func (h *DataHandler) Handle(ctx context.Context, msg z2m.Msg) {
	var sensorMsg doorSensorMsg

	err := json.Unmarshal(msg.Payload, &sensorMsg)
	if err != nil {
		h.l.Error("failed to unmarshal door sensor message", "err", err)
		return
	}

	// contact is true when the door is closed
	if sensorMsg.Contact {
		h.deviceNotifier.SetOpened(ctx, msg.Device, false)
	} else {
		h.deviceNotifier.SetOpened(ctx, msg.Device, true)
	}
}

type doorSensorMsg struct {
	Battery           int  `json:"battery"`
	Contact           bool `json:"contact"`
	DeviceTemperature int  `json:"device_temperature"`
	LinkQuality       int  `json:"linkquality"`
	PowerOutageCount  int  `json:"power_outage_count"`
	TriggerCount      int  `json:"trigger_count"`
	Voltage           int  `json:"voltage"`
}
