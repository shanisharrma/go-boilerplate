package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/nrredis-v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/shanisharrma/go-boilerplate/internal/app/core/repository"
	"github.com/shanisharrma/go-boilerplate/internal/app/core/service"
	"github.com/shanisharrma/go-boilerplate/internal/app/http/handler"
	"github.com/shanisharrma/go-boilerplate/internal/app/http/router"
	"github.com/shanisharrma/go-boilerplate/internal/app/worker/job"
	"github.com/shanisharrma/go-boilerplate/internal/infrastructure/database"
	loggerPkg "github.com/shanisharrma/go-boilerplate/internal/infrastructure/logger"
	"github.com/shanisharrma/go-boilerplate/internal/server"
	"github.com/shanisharrma/go-boilerplate/internal/shared/config"
)

type Application struct {
	Server *server.Server
	Logger zerolog.Logger
}

func Bootstrap(ctx context.Context) (*Application, error) {
	// 1. Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %s", err.Error())
	}

	// 2. Logger / Observability
	loggerService := loggerPkg.NewLoggerService(cfg.Observability)
	// (optional but recommended)
	defer loggerService.Shutdown()

	log := loggerPkg.NewLoggerWithService(cfg.Observability, loggerService)
	// 3. Database migration
	if cfg.Primary.Env != "local" {
		if err := database.Migrate(ctx, &log, cfg); err != nil {
			log.Fatal().Err(err).Msg("failed to migrate database!")
			return nil, err
		}
	}

	// 4. Database
	db, err := database.New(cfg, &log, loggerService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// 5. Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
	})

	if loggerService.GetApplication() != nil {
		redisClient.AddHook(nrredis.NewHook(redisClient.Options()))
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(pingCtx).Err(); err != nil {
		log.Warn().Err(err).Msg("redis unavailable, continuing without redis")
	}

	// 6. Jobs
	jobService := job.NewJobService(&log, cfg)
	jobService.InitHandlers(cfg, &log)

	if err := jobService.Start(); err != nil {
		return nil, err
	}

	// 7. Server (runtime container)
	srv := server.New(
		cfg,
		&log,
		loggerService,
		db,
		redisClient,
		jobService,
	)

	// 8. Repositories
	repos := repository.NewRepositories(srv)

	// 9. Services
	services, err := service.NewServices(srv, repos)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create services")
		return nil, err
	}

	// 10. Handlers
	handlers := handler.NewHandlers(srv, services)

	// 11. Router
	r := router.NewRouter(srv, handlers, services)

	// 12. HTTP server
	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	srv.SetupHTTPServer(httpServer)

	return &Application{
		Server: srv,
		Logger: log,
	}, nil
}
