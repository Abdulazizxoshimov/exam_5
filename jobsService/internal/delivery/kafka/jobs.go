package kafka

import (
	"context"
	"encoding/json"
	"exam_5/jobsService/internal/entity"
	"exam_5/jobsService/internal/infrastructure/kafka"
	"exam_5/jobsService/internal/pkg/config"
	"exam_5/jobsService/internal/usecase"
	"exam_5/jobsService/internal/usecase/event"
	"fmt"

	"github.com/k0kubun/pp"
	"go.uber.org/zap"
)

type jobCreateHandler struct {
	config         *config.Config
	brokerConsumer event.BrokerConsumer
	logger         *zap.Logger
	jobUsecase    usecase.Job
}

func NewJobCreateHandler(config *config.Config,
	brokerConsumer event.BrokerConsumer,
	logger *zap.Logger,
	jobUsecase usecase.Job) *jobCreateHandler {
	return &jobCreateHandler{
		config:         config,
		brokerConsumer: brokerConsumer,
		logger:         logger,
		jobUsecase:    jobUsecase,
	}
}

func (h *jobCreateHandler) HandlerEvents() error {
	consumerConfig := kafka.NewConsumerConfig(
		h.config.Kafka.Address,
		"api.create.user",
		"1",
		func(ctx context.Context, key, value []byte) error {
			var user *entity.Jobs

			if err := json.Unmarshal(value, &user); err != nil {
				return err
			}
			fmt.Println(`
			
			
			kafka message
			
			
			`)
            fmt.Println(user)
			if _, err := h.jobUsecase.Create(ctx, user); err != nil {
				pp.Println("helloodododldldkdkdlsldkfkd")
				return err
			}

			return nil
		},
	)

	h.brokerConsumer.RegisterConsumer(consumerConfig)
	h.brokerConsumer.Run()

	return nil

}
