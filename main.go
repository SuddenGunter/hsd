package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/SuddenGunter/hsd/app"
	"github.com/SuddenGunter/hsd/app/config"
	"github.com/SuddenGunter/hsd/exitcode"
)

var verbose = flag.Bool("v", false, "enable verbose logs")

func main() {
	flag.Parse()

	l := slog.New(slog.NewTextHandler(os.Stdout, logOpts()))

	cfg, err := config.LoadEnv()
	if err != nil {
		l.Error("failed to load config", "err", err)
		os.Exit(exitcode.LoadConfig)
	}

	l.Debug("config loaded", "cfg", cfg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.New(l, cfg).Run(ctx)
}

func logOpts() *slog.HandlerOptions {
	if *verbose {
		return &slog.HandlerOptions{Level: slog.LevelDebug}
	}

	return nil
}
