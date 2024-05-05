package app

import (
	"fmt"
	k "exam_5/jobsService/internal/delivery/kafka"
	"exam_5/jobsService/internal/infrastructure/kafka"
	PostgresDB "exam_5/jobsService/internal/infrastructure/repository/postgres"
	"exam_5/jobsService/internal/pkg/config"
	logPkg "exam_5/jobsService/internal/pkg/logger"
	"exam_5/jobsService/internal/pkg/postgres"
	"exam_5/jobsService/internal/usecase"
	"exam_5/jobsService/internal/usecase/event"

	"go.uber.org/zap"
)

type JobCreateConsumerCLI struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *postgres.PostgresDB
	BrokerConsumer event.BrokerConsumer
}

func NewJobCreateConsumer(config *config.Config) (*JobCreateConsumerCLI, error) {
	logger, err := logPkg.New(config.LogLevel, config.Environment, config.APP+"_cli"+".log")
	if err != nil {
		return nil, err
	}

	consumer := kafka.NewConsumer(logger)

	db, err := postgres.New(config)
	if err != nil {
		return nil, err
	}

	return &JobCreateConsumerCLI{
		Config:         config,
		DB:             db,
		Logger:         logger,
		BrokerConsumer: consumer,
	}, nil
}

func (c *JobCreateConsumerCLI) Run() error {
	fmt.Print("consume is running ....")
	// repo init
	userRepo := PostgresDB.NewJobRepo(c.DB)

	// usecase init
	userUsecase := usecase.NewJobService(c.DB.Config().ConnConfig.ConnectTimeout, userRepo)

	eventHandler := k.NewJobCreateHandler(c.Config, c.BrokerConsumer, c.Logger, userUsecase)

	return eventHandler.HandlerEvents()
}

func (c *JobCreateConsumerCLI) Close() {
	c.BrokerConsumer.Close()

	c.Logger.Sync()
}
