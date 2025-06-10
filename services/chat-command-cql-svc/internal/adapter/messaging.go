package adapter

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/nats-io/nats.go"
)

type Messaging interface {
	AuthRequired() bool
	Barrier(f func()) error
	Buffered() (int, error)
	ChanQueueSubscribe(subj string, group string, ch chan *nats.Msg) (*nats.Subscription, error)
	ChanSubscribe(subj string, ch chan *nats.Msg) (*nats.Subscription, error)
	Close()
	ClosedHandler() nats.ConnHandler
	ConnectedAddr() string
	ConnectedClusterName() string
	ConnectedServerId() string
	ConnectedServerName() string
	ConnectedServerVersion() string
	ConnectedUrl() string
	ConnectedUrlRedacted() string
	DisconnectErrHandler() nats.ConnErrHandler
	DiscoveredServers() []string
	DiscoveredServersHandler() nats.ConnHandler
	Drain() error
	ErrorHandler() nats.ErrHandler
	Flush() error
	FlushTimeout(timeout time.Duration) (err error)
	FlushWithContext(ctx context.Context) error
	ForceReconnect() error
	GetClientID() (uint64, error)
	GetClientIP() (net.IP, error)
	HeadersSupported() bool
	IsClosed() bool
	IsConnected() bool
	IsDraining() bool
	IsReconnecting() bool
	JetStream(opts ...nats.JSOpt) (nats.JetStreamContext, error)
	LastError() error
	LocalAddr() string
	MaxPayload() int64
	NewInbox() string
	NewRespInbox() string
	NumSubscriptions() int
	Publish(subj string, data []byte) error
	PublishMsg(m *nats.Msg) error
	PublishRequest(subj string, reply string, data []byte) error
	QueueSubscribe(subj string, queue string, cb nats.MsgHandler) (*nats.Subscription, error)
	QueueSubscribeSync(subj string, queue string) (*nats.Subscription, error)
	QueueSubscribeSyncWithChan(subj string, queue string, ch chan *nats.Msg) (*nats.Subscription, error)
	RTT() (time.Duration, error)
	ReconnectHandler() nats.ConnHandler
	RemoveStatusListener(ch chan nats.Status)
	Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error)
	RequestMsg(msg *nats.Msg, timeout time.Duration) (*nats.Msg, error)
	RequestMsgWithContext(ctx context.Context, msg *nats.Msg) (*nats.Msg, error)
	RequestWithContext(ctx context.Context, subj string, data []byte) (*nats.Msg, error)
	Servers() []string
	SetClosedHandler(cb nats.ConnHandler)
	SetDisconnectErrHandler(dcb nats.ConnErrHandler)
	SetDisconnectHandler(dcb nats.ConnHandler)
	SetDiscoveredServersHandler(dscb nats.ConnHandler)
	SetErrorHandler(cb nats.ErrHandler)
	SetReconnectHandler(rcb nats.ConnHandler)
	Stats() nats.Statistics
	Status() nats.Status
	StatusChanged(statuses ...nats.Status) chan nats.Status
	Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error)
	SubscribeSync(subj string) (*nats.Subscription, error)
	TLSConnectionState() (tls.ConnectionState, error)
	TLSRequired() bool
}
