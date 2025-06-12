package adapter

import (
	"context"
	"errors"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
)

type MessagingAdapter interface {
	PublishMessage(ctx context.Context, topic string, data any) error
	Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error)
}

type messagingAdapter struct {
	nats *nats.Conn
}

func NewMessagingAdapter(nats *nats.Conn) MessagingAdapter {
	return &messagingAdapter{nats: nats}
}

func (a *messagingAdapter) PublishMessage(ctx context.Context, topic string, data any) error {
	payload, err := sonic.ConfigFastest.Marshal(data)
	if err != nil {
		return errors.New("failed to marshal data: " + err.Error())
	}

	if err := a.nats.Publish(topic, payload); err != nil {
		return errors.New("failed publish message using nats: " + err.Error())
	}

	return nil
}

func (a *messagingAdapter) Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	sub, err := a.nats.Subscribe(subj, cb)
	if err != nil {
		return nil, errors.New("failed to subscribe to subject: " + err.Error())
	}

	return sub, nil
}
