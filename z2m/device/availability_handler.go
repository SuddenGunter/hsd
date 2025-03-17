package device

import (
	"context"
	"log/slog"

	"github.com/SuddenGunter/hsd/z2m"
)

type deviceNotifier interface {
	SetAvailability(ctx context.Context, device string, available bool)
	SetOpened(ctx context.Context, device string, opened bool)
}

// AvailabilityHandler handles availability messages from zigbee2mqtt.
// https://www.zigbee2mqtt.io/guide/configuration/device-availability.html
// It supports both v1 and v2 availability payloads.
type AvailabilityHandler struct {
	deviceNotifier deviceNotifier
	l              *slog.Logger
}

// NewAvailabilityHandler returns a new AvailabilityHandler.
func NewAvailabilityHandler(deviceNotifier deviceNotifier, l *slog.Logger) *AvailabilityHandler {
	return &AvailabilityHandler{deviceNotifier: deviceNotifier, l: l}
}

// Handle handles availability messages from zigbee2mqtt.
func (h *AvailabilityHandler) Handle(ctx context.Context, msg z2m.Msg) {
	switch string(msg.Payload) {
	case `{"state":"online"}`, "online":
		h.deviceNotifier.SetAvailability(ctx, msg.Device, true)
	case `{"state":"offline"}`, "offline":
		h.deviceNotifier.SetAvailability(ctx, msg.Device, false)
	default:
		h.l.Error("failed to parse device availability", "device", msg.Device, "payload", string(msg.Payload))
	}
}
