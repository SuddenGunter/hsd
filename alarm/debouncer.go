package alarm

import (
	"log/slog"
	"sync"
	"time"
)

type debouncer struct {
	next alarmer
	l    *slog.Logger

	mux              *sync.Mutex
	last             time.Time
	debounceInterval time.Duration
}

func newDebouncer(next alarmer, l *slog.Logger) *debouncer {
	return &debouncer{
		next:             next,
		l:                l,
		mux:              &sync.Mutex{},
		debounceInterval: 1 * time.Second,
	}
}

func (d *debouncer) Alarm(device, message string) {
	if !d.allowSend() {
		d.l.Info("alarm event received, but got debounced", "device", device)
		return
	}

	d.next.Alarm(device, message)
}

func (d *debouncer) allowSend() bool {
	d.mux.Lock()
	defer d.mux.Unlock()

	if time.Since(d.last) >= d.debounceInterval {
		d.last = time.Now()
		return true
	}

	return false
}
