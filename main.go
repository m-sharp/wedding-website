package main

import (
	"context"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/m-sharp/wedding-website/lib"
	"github.com/m-sharp/wedding-website/lib/migrations"
	"github.com/m-sharp/wedding-website/web"
)

var (
	// Rendering context variables
	pageContext = &web.RenderContext{
		TargetYear: 2023,
		TargetDate: "10.7.23",
	}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := getCfg()
	logger := getLogger(cfg)

	client := getDBClient(ctx, cfg, logger)
	if err := migrations.RunAll(ctx, client, logger); err != nil {
		log.Fatal("Failed to run DB migrations", zap.Error(err))
	}

	server := web.NewWebServer(logger, pageContext)
	if err := server.Serve(); err != nil {
		logger.Panic("Server stopped listening", zap.Error(err))
	}
}

func getCfg() *lib.Config {
	cfg, err := lib.NewConfig()
	if err != nil {
		log.Fatalf("Error creating Config: %s", err.Error())
	}

	return cfg
}

func getLogger(cfg *lib.Config) *zap.Logger {
	dev, err := cfg.Get(lib.Development)
	if err != nil {
		dev = "false"
	}

	var logger *zap.Logger
	if dev == "true" {
		logger, err = zap.NewDevelopment()

	} else {
		conf := zap.NewProductionConfig()
		conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logger, err = conf.Build()
	}
	if err != nil {
		log.Fatalf("Error creating Logger: %s", err.Error())
	}

	return logger
}

func getDBClient(ctx context.Context, cfg *lib.Config, log *zap.Logger) *lib.DBClient {
	client, err := lib.NewDBClient(ctx, cfg, log)
	if err != nil {
		log.Fatal("Error creating DB client", zap.Error(err))
	}
	if err = client.CheckConnection(); err != nil {
		log.Fatal("DB connection check failed", zap.Error(err))
	}

	return client
}
