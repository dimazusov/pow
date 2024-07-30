package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimazusov/pow/internal/config"
	"github.com/dimazusov/pow/internal/server"
)

func main() {
	cfg, err := config.NewServerCfg()
	if err != nil {
		slog.Error("failed to create config", "error", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-ctx.Done()
		cancel()
	}()

	var lvl slog.LevelVar

	if err = lvl.UnmarshalText([]byte(cfg.LoggerLvl)); err != nil {
		slog.Error("failed to create config", "error", err)
		os.Exit(1)
	}

	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl.Level(),
	}))

	slog.SetDefault(l)

	srv, err := server.New(cfg)
	if err != nil {
		slog.Error("failed creating server", "error", err)
		os.Exit(1)
	}

	if err = srv.Start(ctx); err != nil {
		slog.Error("failed starting server", "error", err)
		os.Exit(1)
	}
}
