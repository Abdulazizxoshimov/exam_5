package event

import (
	"context"
	// "exam_5/clientService/internal/entity"
)

type ConsumerConfig interface {
	GetBrokers() []string
	GetTopic() string
	GetGroupID() string
	GetHandler() func(ctx context.Context, key, value []byte) error
}

type BrokerConsumer interface {
	Run()
	RegisterConsumer(config ConsumerConfig)
	Close()
}


type BrokerProducer interface {
	ProduceContent(ctx context.Context, key string, /*value *entity.User*/) error
	Close()
}