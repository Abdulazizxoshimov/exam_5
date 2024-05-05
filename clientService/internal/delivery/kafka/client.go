package kafka

import (
	"context"
	"encoding/json"
	"exam_5/clientService/internal/entity"
	"exam_5/clientService/internal/infrastructure/kafka"
	"exam_5/clientService/internal/pkg/config"
	"exam_5/clientService/internal/usecase"
	"exam_5/clientService/internal/usecase/event"
	"fmt"

	"github.com/k0kubun/pp"
	"go.uber.org/zap"
)

type userCreateHandler struct {
	config         *config.Config
	brokerConsumer event.BrokerConsumer
	logger         *zap.Logger
	userUsecase    usecase.User
}

func NewUserCreateHandler(config *config.Config,
	brokerConsumer event.BrokerConsumer,
	logger *zap.Logger,
	userUsecase usecase.User) *userCreateHandler {
	return &userCreateHandler{
		config:         config,
		brokerConsumer: brokerConsumer,
		logger:         logger,
		userUsecase:    userUsecase,
	}
}

func (h *userCreateHandler) HandlerEvents() error {
	consumerConfig := kafka.NewConsumerConfig(
		h.config.Kafka.Address,
		"api.create.user",
		"1",
		func(ctx context.Context, key, value []byte) error {
			var user *entity.Client

			if err := json.Unmarshal(value, &user); err != nil {
				return err
			}
			fmt.Println(`
			
			
			kafka message
			
			
			`)
            fmt.Println(user)
			if _, err := h.userUsecase.Create(ctx, user); err != nil {
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
