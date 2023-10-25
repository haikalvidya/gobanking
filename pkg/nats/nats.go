package nats

import (
	"gobanking/pkg/logger"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	connectWait   = time.Second * 30
	interval      = time.Second * 10
	maxReconnects = 10
	reconnectWait = 5
	TimeoutReq    = time.Second * 5
)

type NatsConfig struct {
	URL string `yaml:"url"`
}

func NewNatsConnect(cfg *NatsConfig, log logger.Logger) (*nats.Conn, error) {
	nc, err := nats.Connect(
		cfg.URL,
		nats.Timeout(connectWait),
		nats.PingInterval(interval),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(maxReconnects),
		nats.ReconnectWait(reconnectWait),
		nats.ConnectHandler(func(nc *nats.Conn) {
			log.Info("nats: client connected")
		}),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Fatalf("nats: Connection lost: %v", err)
		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			log.Warn("nats: client reconnected")
		}),
		nats.ClosedHandler(func(c *nats.Conn) {
			log.Warn("nats: client connection closed")
		}),
	)

	return nc, err
}
