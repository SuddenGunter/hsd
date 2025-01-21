package alarmposthandler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/SuddenGunter/hsd/alarm"
)

// PostHandler handles POST requests to /alarm.
type PostHandler struct {
	l       *slog.Logger
	alarmer *alarm.Alarmer
}

// NewPostHandler returns a new PostHandler.
func NewPostHandler(l *slog.Logger, alarmer *alarm.Alarmer) *PostHandler {
	return &PostHandler{l, alarmer}
}

// ServeHTTP handles the request.
func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.l.Error("failed to read request body", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		h.l.Error("failed to unmarshal request", "err", err)
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	if req.Enabled {
		h.alarmer.Enable()
	} else {
		h.alarmer.Disable()
	}

	resp, err := json.Marshal(map[string]bool{"enabled": req.Enabled})
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
