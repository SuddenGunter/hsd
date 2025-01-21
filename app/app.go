package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/SuddenGunter/hsd/alarm"
	alarmgethandler "github.com/SuddenGunter/hsd/api/alarm/get"
	alarmposthandler "github.com/SuddenGunter/hsd/api/alarm/post"
	"github.com/SuddenGunter/hsd/app/config"
	"github.com/SuddenGunter/hsd/telegram"
	"github.com/SuddenGunter/hsd/z2m"
	"github.com/SuddenGunter/hsd/z2m/device"
	"github.com/SuddenGunter/hsd/z2m/mqttc"
)

// App manages app/service lifecycle.
// It is responsible for starting and stopping all the dependencies.
type App struct {
	l   *slog.Logger
	cfg *config.Config
}

// New returns a new App.
func New(l *slog.Logger, cfg *config.Config) *App {
	return &App{l, cfg}
}

// Run starts the app and blocks until shutdown.
func (app *App) Run(sigCtx context.Context) {
	notifier := telegram.NewNotifier(app.cfg.TelegramBotToken, app.l)
	alarmer := alarm.New(notifier, app.l)
	devMsg := alarm.NewDeviceMessenger(app.cfg.Z2MDevices, alarmer, app.l)
	devMsg.Listen()
	go notifier.Listen()
	defer devMsg.Close()

	gh := alarmgethandler.NewGetHandler(app.l, alarmer)
	ph := alarmposthandler.NewPostHandler(app.l, alarmer)

	mux := http.NewServeMux()
	mux.Handle("GET /alarm", gh)
	mux.Handle("POST /alarm", ph)

	app.l.Debug("connecting to mqtt broker")

	mc, err := mqttc.Connect(app.cfg)
	if err != nil {
		app.l.Error("failed to connect to mqtt broker", "err", err)
		return
	}

	//nolint:gosec // false positive - this code is only used on 64-bit systems
	defer mc.Disconnect(uint((5 * time.Second).Milliseconds()))

	z2ml := z2m.NewZigbee2MQTTListener(mc, device.NewDataHandler(devMsg, app.l), device.NewAvailabilityHandler(devMsg, app.l), app.cfg.Z2MDevices, app.l)
	z2ml.Subscribe()

	ctx, crash := context.WithCancel(sigCtx)
	srv := http.Server{
		ReadTimeout: 5 * time.Second,
		Addr:        fmt.Sprintf(":%d", app.cfg.Port),
		Handler:     mux,
	}

	go func() {
		app.l.Debug("starting http server", "addr", srv.Addr)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.l.Error("failed to listen and serve", "err", err)
			crash()
		}
	}()

	app.l.Info("app started")
	<-ctx.Done()

	app.l.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		app.l.Info("server shutdown returned an err: %v\n", "err", err)
	}

	app.l.Info("shutdown complete")
}
