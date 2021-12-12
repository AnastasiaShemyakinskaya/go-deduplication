package api

import (
	"fmt"
	"go-deduplication/config"
	"go-deduplication/internal/helper"
	"go-deduplication/internal/services"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Service struct {
	loader    services.Loader
	reader    services.Reader
	restorer  services.Restorer
	allReader helper.AllReader
	cfg       *config.Config
	logger    *zap.Logger
}

func NewService(
	loader services.Loader,
	reader services.Reader,
	restorer services.Restorer,
	allReader helper.AllReader,
	logger *zap.Logger,
	cfg *config.Config,
) *Service {
	return &Service{
		loader:    loader,
		reader:    reader,
		restorer:  restorer,
		allReader: allReader,
		cfg:       cfg,
		logger:    logger,
	}
}

const (
	FileParam = "file"
)

func (srv *Service) loadFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		context := r.Context()
		file := query.Get(FileParam)
		fileContent, err := srv.allReader.ReadAll(file)
		if err != nil {
			srv.logger.Error("failed to read whole file", zap.Error(err), zap.String("filename", file))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		now := time.Now()
		err = srv.loader.LoadFile(context, srv.cfg.FileDirectory, file, srv.cfg.HashFunction, srv.cfg.ByteSize, fileContent)
		if err != nil {
			srv.logger.Error("failed to load file in db", zap.Error(err), zap.String("filename", file))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		srv.logger.Info(fmt.Sprintf("Load time: %s", time.Now().Sub(now).String()))
	})
}

func (srv *Service) readFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		context := r.Context()
		file := query.Get(FileParam)
		data, err := srv.reader.ReadFile(context, srv.cfg.FileDirectory, file, srv.cfg.HashFunction, srv.cfg.ByteSize)
		if err != nil {
			srv.logger.Error("failed to read file", zap.Error(err), zap.String("filename", file))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
}

func (srv *Service) restoreFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		context := r.Context()
		file := query.Get(FileParam)
		fileContent, err := srv.allReader.ReadAll(file)
		if err != nil {
			srv.logger.Error("failed to read whole file", zap.Error(err), zap.String("filename", file))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		now := time.Now()
		_, err = srv.restorer.RestoreFile(context, fileContent, file, srv.cfg.FileDirectory, srv.cfg.HashFunction, srv.cfg.ByteSize)
		if err != nil {
			srv.logger.Error("failed to restore file", zap.Error(err), zap.String("filename", file))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		srv.logger.Info(fmt.Sprintf("Load time: %s", time.Now().Sub(now).String()))
	})
}

func (srv *Service) PrepareHandler(r *http.ServeMux) {
	r.Handle("/v1/load", srv.loadFile())
	r.Handle("/v1/read", srv.readFile())
	r.Handle("/v1/restore", srv.restoreFile())
}
