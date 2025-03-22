package alarm

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type alarmer interface {
	Alarm(device, message string)
}

// Device that can be alarmed. Processes state updates and alarms if necessary.
// For not the only type of a supported device is a door sensor.
type Device struct {
	alarmer alarmer

	name        string
	available   bool
	opened      bool
	lastUpdated int64

	stateUpdate chan stateUpdateMsg
	close       chan struct{}

	l *slog.Logger
}

type stateUpdateMsg struct {
	availability *bool
	opened       *bool
}

// NewDevice returns a new Device.
func NewDevice(name string, alarmer alarmer, l *slog.Logger) *Device {
	return &Device{
		name:    name,
		alarmer: alarmer,
		// we assume it's available unless we hear otherwise
		available:   true,
		opened:      false,
		stateUpdate: make(chan stateUpdateMsg),
		close:       make(chan struct{}),
		l:           l,
	}
}

// SetAvailability sets the availability of the device.
func (d *Device) SetAvailability(ctx context.Context, available bool) {
	select {
	case <-ctx.Done():
		d.l.Error("device state update timeout", "device", d.name, "operation", "SetAvailability", "reason", ctx.Err())
		return
	case d.stateUpdate <- stateUpdateMsg{availability: &available}:
		return
	case <-d.close:
		return
	}
}

// SetOpened sets the opened state of the device.
func (d *Device) SetOpened(ctx context.Context, opened bool) {
	select {
	case <-ctx.Done():
		d.l.Error("device state update timeout", "device", d.name, "operation", "SetOpened", "reason", ctx.Err())
		return
	case d.stateUpdate <- stateUpdateMsg{opened: &opened}:
		return
	case <-d.close:
		return
	}
}

func (d *Device) loop() {
	for {
		select {
		case <-d.close:
			return

		case msg := <-d.stateUpdate:
			d.l.Info("device state update received", "device", d.name, "availability", ptr(msg.availability), "opened", ptr(msg.opened))

			if msg.availability != nil {
				d.available = *msg.availability
			}

			if msg.opened != nil {
				d.opened = *msg.opened
			}

			d.lastUpdated = time.Now().Unix()

			d.evalAlarm()
		}
	}
}

func (d *Device) evalAlarm() {
	if d.opened {
		d.alarmer.Alarm(d.name, "opened")
		return
	}

	if !d.available {
		d.alarmer.Alarm(d.name, "unavailable")
		return
	}

	if d.lastUpdated < time.Now().Add(-time.Hour*26).Unix() {
		d.alarmer.Alarm(d.name, "no messages received for a long time")
		return
	}
}

func ptr(b *bool) string {
	if b == nil {
		return "nil"
	}

	return fmt.Sprintf("%v", *b)
}
