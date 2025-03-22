package alarm

import (
	"context"
	"log/slog"
)

// DeviceMessenger is a collection of devices that can be alarmed.
// Handles creation of devices, their state updates and lifecycle.
type DeviceMessenger struct {
	devices map[string]*Device

	l *slog.Logger
}

// NewDeviceMessenger returns a new DeviceMessenger.
func NewDeviceMessenger(devices []string, alarmer alarmer, l *slog.Logger) *DeviceMessenger {
	d := make(map[string]*Device)
	for _, device := range devices {
		d[device] = NewDevice(device, NewDebouncer(alarmer, l), l)
	}

	return &DeviceMessenger{devices: d, l: l}
}

// SetAvailability sets the availability of the device.
func (m *DeviceMessenger) SetAvailability(ctx context.Context, device string, available bool) {
	if d, ok := m.devices[device]; ok {
		d.SetAvailability(ctx, available)
	} else {
		m.l.Error("device not found", "device", device, "operation", "SetAvailability")
	}
}

// SetOpened sets the opened state of the device.
func (m *DeviceMessenger) SetOpened(ctx context.Context, device string, opened bool) {
	if d, ok := m.devices[device]; ok {
		d.SetOpened(ctx, opened)
	} else {
		m.l.Error("device not found", "device", device, "operation", "SetOpened")
	}
}

// Close closes all devices.
func (m *DeviceMessenger) Close() {
	for _, d := range m.devices {
		close(d.close)
	}
}

// Listen starts listening for device state updates.
func (m *DeviceMessenger) Listen() {
	for _, d := range m.devices {
		go d.loop()
	}
}
