package main

import (
	"github.com/fadyat/erida/internal"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
)

func readConfig() (*internal.Config, error) {
	var cfg internal.Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func main() {
	cfg, err := readConfig()
	if err != nil {
		slog.Error("failed to read config", "error", err)
		return
	}

	srv := internal.NewServer(cfg)
	slog.Info("Starting server at: ", "addr", srv.Addr)
	if err = srv.ListenAndServe(); err != nil {
		slog.Error("failed to start server: ", err)
	}
}
