package alarm

import (
	"log/slog"
	"sync"
	"time"
)

type Debouncer struct {
	next alarmer
	l    *slog.Logger

	mux              *sync.Mutex
	last             time.Time
	debounceInterval time.Duration
}

func NewDebouncer(next alarmer, l *slog.Logger) *Debouncer {
	return &Debouncer{
		next:             next,
		l:                l,
		mux:              &sync.Mutex{},
		debounceInterval: 1 * time.Second,
	}
}

func (d *Debouncer) Alarm(device, message string) {
	if !d.allowSend() {
		d.l.Info("alarm event received, but got debounced", "device", device)
		return
	}

	d.next.Alarm(device, message)
}

func (d *Debouncer) allowSend() bool {
	d.mux.Lock()
	defer d.mux.Unlock()

	if time.Since(d.last) >= d.debounceInterval {
		d.last = time.Now()
		return true
	}

	return false
}
