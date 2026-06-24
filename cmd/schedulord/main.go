package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/k3nsei/schedulord/internal/config"
	"github.com/k3nsei/schedulord/internal/scheduler"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	token, err := config.ReadGithubToken(cfg)
	if err != nil {
		log.Fatal(err)
	}

	loc, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		loc = time.UTC
	}

	s := scheduler.New(loc, token)

	if err := s.Start(cfg); err != nil {
		log.Fatal(err)
	}

	slog.Info("Schedulord is running")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-stop
}
