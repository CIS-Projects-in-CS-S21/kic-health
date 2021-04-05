package main

import (
	"context"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/kic/health/internal/setup"
	"github.com/kic/health/pkg/logging"
)

func main() {
	IsProduction := os.Getenv("PRODUCTION") != ""
	var logger *zap.SugaredLogger
	if IsProduction {
		logger = logging.CreateLogger(zapcore.InfoLevel)
	} else {
		logger = logging.CreateLogger(zapcore.DebugLevel)
	}


	repo, mongoClient := setup.DBRepositorySetup(logger, "health")

	serv := setup.GRPCSetup(logger, repo)

	defer serv.Stop()
	defer mongoClient.Disconnect(context.Background())

	// the server is listening in a goroutine so hang until we get an interrupt signal
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
}