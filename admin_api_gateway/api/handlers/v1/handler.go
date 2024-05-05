package v1

import (
	grpcClients "exam_5/admin_api_gateway/internal/infrastructure/grpc_service_client"
	"exam_5/admin_api_gateway/internal/pkg/config"
	repo "exam_5/admin_api_gateway/internal/infrastructure/repository/redis"
	t "exam_5/admin_api_gateway/internal/pkg/token"
	"exam_5/admin_api_gateway/internal/usecase/event"
	"time"

	"go.uber.org/zap"
)

type HandlerV1 struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Service        grpcClients.ServiceClient
	BrokerProducer event.BrokerProducer
	redisStorage   repo.Cache
	jwthandler     t.JWTHandler
}

// HandlerV1Config ...
type HandlerV1Config struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Service        grpcClients.ServiceClient
	BrokerProducer event.BrokerProducer
	Redis          repo.Cache
	JWTHandler     t.JWTHandler
}

// New ...
func New(c *HandlerV1Config) *HandlerV1 {
	return &HandlerV1{
		Config:         c.Config,
		Logger:         c.Logger,
		Service:        c.Service,
		ContextTimeout: c.ContextTimeout,
		BrokerProducer: c.BrokerProducer,
		redisStorage:   c.Redis,
		jwthandler:     c.JWTHandler,
	}
}
