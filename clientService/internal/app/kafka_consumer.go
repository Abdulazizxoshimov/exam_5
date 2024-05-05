package app

import (
	"fmt"
	k "exam_5/clientService/internal/delivery/kafka"
	"exam_5/clientService/internal/infrastructure/kafka"
	PostgresDB "exam_5/clientService/internal/infrastructure/repository/postgres"
	"exam_5/clientService/internal/pkg/config"
	logPkg "exam_5/clientService/internal/pkg/logger"
	"exam_5/clientService/internal/pkg/postgres"
	"exam_5/clientService/internal/usecase"
	"exam_5/clientService/internal/usecase/event"

	"go.uber.org/zap"
)

type UserCreateConsumerCLI struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *postgres.PostgresDB
	BrokerConsumer event.BrokerConsumer
}

func NewUserCreateConsumer(config *config.Config) (*UserCreateConsumerCLI, error) {
	logger, err := logPkg.New(config.LogLevel, config.Environment, config.APP+"_cli"+".log")
	if err != nil {
		return nil, err
	}

	consumer := kafka.NewConsumer(logger)

	db, err := postgres.New(config)
	if err != nil {
		return nil, err
	}

	return &UserCreateConsumerCLI{
		Config:         config,
		DB:             db,
		Logger:         logger,
		BrokerConsumer: consumer,
	}, nil
}

func (c *UserCreateConsumerCLI) Run() error {
	fmt.Print("consume is running ....")
	// repo init
	userRepo := PostgresDB.NewUserRepo(c.DB)

	// usecase init
	userUsecase := usecase.NewUserService(c.DB.Config().ConnConfig.ConnectTimeout, userRepo)

	eventHandler := k.NewUserCreateHandler(c.Config, c.BrokerConsumer, c.Logger, userUsecase)

	return eventHandler.HandlerEvents()
}

func (c *UserCreateConsumerCLI) Close() {
	c.BrokerConsumer.Close()

	c.Logger.Sync()
}
