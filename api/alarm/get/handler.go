package alarmgethandler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/SuddenGunter/hsd/alarm"
)

// GetHandler handles GET requests to /alarm.
type GetHandler struct {
	l       *slog.Logger
	alarmer *alarm.Alarmer
}

// NewGetHandler returns a new GetHandler.
func NewGetHandler(l *slog.Logger, alarmer *alarm.Alarmer) *GetHandler {
	return &GetHandler{l, alarmer}
}

// ServeHTTP handles the request.
func (h *GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := h.alarmer.Enabled()

	resp, err := json.Marshal(map[string]bool{"enabled": v})
	if err != nil {
		h.l.Error("failed to marshal response", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	_, err = w.Write(resp)
	if err != nil {
		h.l.Warn("failed to write response", "err", err, "path", r.URL.Path)
	}
}
