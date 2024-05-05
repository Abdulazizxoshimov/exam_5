package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	// "github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"exam_5/admin_api_gateway/api"

	grpcService "exam_5/admin_api_gateway/internal/infrastructure/grpc_service_client"
	"exam_5/admin_api_gateway/internal/infrastructure/kafka"
	"exam_5/admin_api_gateway/internal/usecase/event"

	// "exam_5/admin_api_gateway/internal/infrastructure/kafka"
	// "exam_5/admin_api_gateway/internal/infrastructure/repository/postgresql"
	redisrepo "exam_5/admin_api_gateway/internal/infrastructure/repository/redis"
	"exam_5/admin_api_gateway/internal/pkg/config"
	"exam_5/admin_api_gateway/internal/pkg/logger"

	"exam_5/admin_api_gateway/internal/pkg/otlp"
	// "exam_5/admin_api_gateway/internal/pkg/policy"
	// "exam_5/admin_api_gateway/internal/pkg/postgres"
	"exam_5/admin_api_gateway/internal/pkg/storage/redis"
	// "exam_5/admin_api_gateway/internal/usecase/app_version"
	// "exam_5/admin_api_gateway/internal/usecase/event"
	// "evrone_api_gateway/internal/usecase/refresh_token"
)

type App struct {
	Config         *config.Config
	Logger         *zap.Logger
	server         *http.Server
	RedisDB        *redis.RedisDB
    ShutdownOTLP   func() error
	Clients        grpcService.ServiceClient
	BrokerProducer event.BrokerProducer
}

func NewApp(cfg config.Config) (*App, error) {
	// logger init
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	// kafka producer init
	kafkaProducer := kafka.NewProducer(&cfg, logger)
	// redis init
	redisdb, err := redis.New(&cfg)
	if err != nil {
		return nil, err
	}
	

	// otlp collector init
	shutdownOTLP, err := otlp.InitOTLPProvider(&cfg)
	if err != nil {
		return nil, err
	}


	return &App{
		Config:         &cfg,
		Logger:         logger,
		ShutdownOTLP:   shutdownOTLP,
		BrokerProducer: kafkaProducer,
		RedisDB:        redisdb,

	}, nil
}

func (a *App) Run() error {
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error while parsing context timeout: %v", err)
	}

	clients, err := grpcService.New(a.Config)
	if err != nil {
		return err
	}
	a.Clients = clients

	// initialize cache
	cache := redisrepo.NewCache(a.RedisDB)

	

	// api init
	handler := api.NewRoute(api.RouteOption{
		Config:         a.Config,
		Logger:         a.Logger,
		ContextTimeout: contextTimeout,
		Service:        clients,
	    BrokerProducer: a.BrokerProducer,
		Cache:          cache,

	})
	// if err = a.Enforcer.LoadPolicy(); err != nil {
	// 	return fmt.Errorf("error during enforcer load policy: %w", err)
	// }

	// server init
	a.server, err = api.NewServer(a.Config, handler)
	if err != nil {
		return fmt.Errorf("error while initializing server: %v", err)
	}
	
	return a.server.ListenAndServe()
}

func (a *App) Stop() {

	// close grpc connections
	a.Clients.Close()

	// shutdown server http
	if err := a.server.Shutdown(context.Background()); err != nil {
		a.Logger.Error("shutdown server http ", zap.Error(err))
	}

	// shutdown otlp collector
	if err := a.ShutdownOTLP(); err != nil {
		a.Logger.Error("shutdown otlp collector", zap.Error(err))
	}

	// zap logger sync
	a.Logger.Sync()
}
