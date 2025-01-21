package alarm

import (
	"log/slog"
	"sync/atomic"
)

type notifier interface {
	Notify(device, msg string)
}

// Alarmer tracks if the alarm is enabled.
// Right now it is a simple flag, but may be extended in the future with scheduling or other features.
type Alarmer struct {
	notifier notifier
	enabled  *atomic.Bool

	l *slog.Logger
}

// New returns a new Alarmer.
func New(notifier notifier, l *slog.Logger) *Alarmer {
	// start with alarm enabled on restart
	e := &atomic.Bool{}
	e.Store(true)

	return &Alarmer{
		notifier: notifier,
		enabled:  e,
		l:        l,
	}
}

// Enabled returns the current state of the alarm.
func (a *Alarmer) Enabled() bool {
	return a.enabled.Load()
}

// Enable the alarm.
func (a *Alarmer) Enable() {
	a.enabled.Store(true)
	a.notifier.Notify("alarm", "enabled")
}

// Disable the alarm.
func (a *Alarmer) Disable() {
	a.enabled.Store(false)
	a.notifier.Notify("alarm", "disabled")
}

// Alarm sends an alarm message to the notifier if the alarm is enabled.
func (a *Alarmer) Alarm(device, message string) {
	if a.Enabled() {
		a.notifier.Notify(device, message)
	} else {
		a.l.Debug("alarm event received, but will be ignored: alarm disabled")
	}
}
