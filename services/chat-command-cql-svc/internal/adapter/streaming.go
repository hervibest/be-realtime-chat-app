package adapter

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Streaming interface {
	AbortTransaction(ctx context.Context) error
	BeginTransaction() error
	Close()
	CommitTransaction(ctx context.Context) error
	Events() chan kafka.Event
	Flush(timeoutMs int) int
	GetFatalError() error
	GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error)
	InitTransactions(ctx context.Context) error
	IsClosed() bool
	Len() int
	Logs() chan kafka.LogEvent
	OffsetsForTimes(times []kafka.TopicPartition, timeoutMs int) (offsets []kafka.TopicPartition, err error)
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
	ProduceChannel() chan *kafka.Message
	Purge(flags int) error
	QueryWatermarkOffsets(topic string, partition int32, timeoutMs int) (low int64, high int64, err error)
	SendOffsetsToTransaction(ctx context.Context, offsets []kafka.TopicPartition, consumerMetadata *kafka.ConsumerGroupMetadata) error
	SetOAuthBearerToken(oauthBearerToken kafka.OAuthBearerToken) error
	SetOAuthBearerTokenFailure(errstr string) error
	SetSaslCredentials(username string, password string) error
	String() string
	TestFatalError(code kafka.ErrorCode, str string) kafka.ErrorCode
}
