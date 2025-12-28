package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/shanisharrma/go-boilerplate/internal/infrastructure/database"
	loggerPkg "github.com/shanisharrma/go-boilerplate/internal/infrastructure/logger"
	"github.com/shanisharrma/go-boilerplate/internal/shared/config"
)

type JobRunner interface {
	Start() error
	Stop()
}

type Server struct {
	Config        *config.Config
	Logger        *zerolog.Logger
	LoggerService *loggerPkg.LoggerService
	DB            *database.Database
	Redis         *redis.Client
	httpServer    *http.Server
	Jobs          JobRunner
}

func New(
	cfg *config.Config,
	logger *zerolog.Logger,
	loggerService *loggerPkg.LoggerService,
	db *database.Database,
	redis *redis.Client,
	jobs JobRunner,
) *Server {
	return &Server{
		Config:        cfg,
		Logger:        logger,
		LoggerService: loggerService,
		DB:            db,
		Redis:         redis,
		Jobs:          jobs,
	}
}

func (s *Server) SetupHTTPServer(server *http.Server) {
	s.httpServer = server
}

func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("HTTPS server not initialized")
	}

	s.Logger.Info().
		Str("port", s.Config.Server.Port).
		Str("env", s.Config.Primary.Env).
		Msg("starting server")

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown http server: %w", err)
		}
	}

	if s.Jobs != nil {
		s.Jobs.Stop()
	}

	if s.Redis != nil {
		_ = s.Redis.Close()
	}

	if s.DB != nil {
		_ = s.DB.Close()
	}

	return nil
}
