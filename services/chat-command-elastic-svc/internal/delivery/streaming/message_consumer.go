package consumer

import (
	"be-realtime-chat-app/services/chat-command-elastic-svc/internal/helper/logs"
	"be-realtime-chat-app/services/chat-command-elastic-svc/internal/model/event"
	"be-realtime-chat-app/services/chat-command-elastic-svc/internal/usecase"
	"context"

	"github.com/bytedance/sonic"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type MessageConsumer interface {
	Consume(message *kafka.Message) error
}

type messageConsumerImpl struct {
	commandAsyncUseCase usecase.CommandAsyncUseCase
	log                 logs.Log
}

func NewMessageConsumerImpl(commandAsyncUseCase usecase.CommandAsyncUseCase, log logs.Log) *messageConsumerImpl {
	return &messageConsumerImpl{
		commandAsyncUseCase: commandAsyncUseCase,
		log:                 log,
	}
}

func (c messageConsumerImpl) Consume(message *kafka.Message) error {
	MessageEvent := new(event.Message)
	if err := sonic.ConfigFastest.Unmarshal(message.Value, MessageEvent); err != nil {
		c.log.Warn("error unmarshalling Message event", zap.Error(err), zap.String("message", string(message.Value)))
		return err
	}
	c.log.Warn("Received topic messages with event", zap.Any("event", MessageEvent), zap.Int32("partition", message.TopicPartition.Partition))

	if err := c.commandAsyncUseCase.Persist(context.TODO(), MessageEvent); err != nil {
		c.log.Error("error persisting Message event", zap.Error(err), zap.String("message", string(message.Value)))
		return err
	}

	// TODO process event
	c.log.Warn("Received topic messages with event", zap.Any("event", MessageEvent), zap.Int32("partition", message.TopicPartition.Partition))
	return nil
}
