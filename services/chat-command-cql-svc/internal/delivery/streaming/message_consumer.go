package consumer

import (
	"be-realtime-chat-app/services/chat-command-cql-svc/internal/model/event"
	"be-realtime-chat-app/services/chat-command-cql-svc/internal/usecase"
	"be-realtime-chat-app/services/commoner/logs"
	"context"

	"github.com/bytedance/sonic"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type MessageConsumer interface {
	Consume(message *kafka.Message) error
}

type messageConsumerImpl struct {
	commandUseCase usecase.CommandUseCase
	log            logs.Log
}

func NewMessageConsumerImpl(commandUseCase usecase.CommandUseCase, log logs.Log) *messageConsumerImpl {
	return &messageConsumerImpl{
		commandUseCase: commandUseCase,
		log:            log,
	}
}

func (c messageConsumerImpl) Consume(message *kafka.Message) error {
	MessageEvent := new(event.Message)
	if err := sonic.ConfigFastest.Unmarshal(message.Value, MessageEvent); err != nil {
		c.log.Warn("error unmarshalling Message event", zap.Error(err), zap.String("message", string(message.Value)))
		return err
	}

	if err := c.commandUseCase.PersistChat(context.TODO(), MessageEvent); err != nil {
		c.log.Error("error persisting Message event", zap.Error(err), zap.String("message", string(message.Value)))
		return err
	}

	// TODO process event
	c.log.Warn("Received topic messages with event", zap.Any("event", MessageEvent), zap.Int32("partition", message.TopicPartition.Partition))
	return nil
}
