package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sophie-rigg/message-store/cache"
	"github.com/sophie-rigg/message-store/server"
)

var (
	port     int
	logLevel string
)

func init() {
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error, fatal, panic)")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	logger := log.With().Ctx(ctx).Fields(map[string]interface{}{
		"log_level": logLevel,
		"port":      port,
	}).Logger()

	// Set the log level
	l, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		logger.Fatal().Err(err).Msg("parsing log level")
	}

	zerolog.SetGlobalLevel(l)

	localCache, err := cache.NewCache()
	if err != nil {
		logger.Fatal().Err(err).Msg("creating cache")
	}

	// Register the handlers
	router := server.Register(localCache)

	logger.Info().Msg("starting server")
	// Start the server
	if err = http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil {
		logger.Fatal().Err(err).Msg("server error")
	}
}
