package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/necutya/meeting-app/config"
	"github.com/necutya/meeting-app/internal/server/handlers"
	"github.com/necutya/meeting-app/internal/server/http"
	"github.com/necutya/meeting-app/internal/service"
	"github.com/necutya/meeting-app/internal/storage/inmemory"

	log "github.com/sirupsen/logrus"
)

var (
	cfg config.Config

	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
)

func init() {
	var err error

	wg = &sync.WaitGroup{}

	cfg, err = config.Read()
	if err != nil {
		panic(err)
	}

	setupLogger(cfg.LogLevel)

	// prepare main context
	ctx, cancel = context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)

}

func setupLogger(logLevel string) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)
	parsedLevel, err := log.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		parsedLevel = log.DebugLevel
	}
	log.SetLevel(parsedLevel)
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChannel
		stop()
	}()
}

func main() {
	srv := service.New(
		&cfg,
		inmemory.NewRoomRepo(),
		// redisClient,
	)

	httpSrv := http.New(
		&cfg.HTTPConfig,
		handlers.NewMeetingHandler(srv),
	)
	httpSrv.Run(ctx, wg)

	// wait while services work
	wg.Wait()
	log.Info("Service stopped")
}
