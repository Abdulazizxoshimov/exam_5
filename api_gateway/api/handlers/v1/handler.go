package v1

import (
	grpcClients "exam_5/api_gateway/internal/infrastructure/grpc_service_client"
	repo "exam_5/api_gateway/internal/infrastructure/repository/redis"
	"exam_5/api_gateway/internal/pkg/config"
	"exam_5/api_gateway/internal/usecase/event"
	"exam_5/api_gateway/internal/usecase/refresh_token"
	"time"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
)

type HandlerV1 struct {
	Config         config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Service        grpcClients.ServiceClient
	BrokerProducer event.BrokerProducer
	redisStorage   repo.Cache
	RefreshToken   refresh_token.JWTHandler
	Enforcer       *casbin.Enforcer
}

// HandlerV1Config ...
type HandlerV1Config struct {
	Config         config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Service        grpcClients.ServiceClient
	BrokerProducer event.BrokerProducer
	Redis          repo.Cache
	RefreshToken   refresh_token.JWTHandler
	Enforcer       *casbin.Enforcer
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
		Enforcer:       c.Enforcer,
		RefreshToken:   c.RefreshToken,
	}
}
