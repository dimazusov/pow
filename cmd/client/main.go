package main

import (
	"log/slog"
	"os"

	"github.com/dimazusov/pow/internal/client"
	"github.com/dimazusov/pow/internal/config"
)

func main() {
	cfg, err := config.NewClientCfg()
	if err != nil {
		slog.Error("failed to create config", "error", err)
		os.Exit(1)
	}

	cl := client.New(cfg)
	for i := 0; i < cfg.RequestCount; i++ {
		quote, err := cl.GetRandomQuote()
		if err != nil {
			slog.Error("failed get random quote", "error", err)
			os.Exit(1)
		}

		slog.Info("response from server", "quote", quote)
	}
}
