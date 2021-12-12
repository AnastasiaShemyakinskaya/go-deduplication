package main

import (
	"go-deduplication/api"
	"go-deduplication/config"
	"go-deduplication/internal/helper"
	"go-deduplication/internal/repository"
	"go-deduplication/internal/services"
	"go-deduplication/internal/systems"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	FileName = "config.yml"
	FilePath = "."
)

func main() {
	logger, _ := zap.NewProduction()
	cfg, err := config.InitConfig(FileName, FilePath)
	if err != nil {
		logger.Fatal("failed read config file", zap.Error(err))
	}
	postgres, err := systems.NewDbConn(cfg.Postgres)
	if err != nil {
		logger.Fatal("failed to connect to postgres", zap.Error(err))
	}
	defer postgres.DB.Close()

	fileRepo := repository.NewFileRepo(postgres.DB)
	hashRepo := repository.NewHashRepo(postgres.DB)
	fileHashRepo := repository.NewFileHashRepo(postgres.DB)
	localFile := systems.NewLocalFile()
	loader := services.NewFileLoader(fileHashRepo, fileRepo, hashRepo, localFile, postgres)
	processor := helper.NewFileProcessor(localFile)
	allReader := helper.NewFileAllReader()
	reader := services.NewFileReader(fileHashRepo, processor)
	restorer := services.NewFileRestorer(fileHashRepo, processor, localFile)
	service := api.NewService(loader, reader, restorer, allReader, logger, cfg)

	mux := http.NewServeMux()
	service.PrepareHandler(mux)

	srv := &http.Server{
		IdleTimeout:  20 * time.Minute,
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Minute,
		Addr:         cfg.Address,
		Handler:      mux,
	}
	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatal("failed to serve http", zap.Error(err))
	}
}
