package nats

import (
	"context"
	"gobanking/pkg/logger"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

const (
	workersNum    = 50
	ackWait       = 1 * time.Minute
	maxAckPending = 1000
)

func NewJetstreamConsumer(ctx context.Context, log logger.Logger, js jetstream.Stream, durableName string, subjects []string) (*jetstream.Consumer, error) {
	consumer, err := js.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:        durableName,
		AckWait:        ackWait,
		MaxAckPending:  maxAckPending,
		AckPolicy:      jetstream.AckExplicitPolicy,
		ReplayPolicy:   jetstream.ReplayInstantPolicy,
		MaxDeliver:     1,
		FilterSubjects: subjects,
	})
	if err != nil {
		log.Errorf("failed to create consumer: %+v", err)
		return nil, err
	}

	return &consumer, nil
}

func ConsumeSubscribePull(
	ctx context.Context,
	consumer jetstream.Consumer,
	handlers map[string]func(ctx context.Context, msg jetstream.Msg) error,
	log logger.Logger,
) {
	iter, err := consumer.Messages(
		jetstream.PullMaxMessages(1),
		jetstream.WithMessagesErrOnMissingHeartbeat(true),
	)
	if err != nil {
		log.Errorf("failed to create iterator: %+v", err)
		return
	}

	defer iter.Stop()

	sem := make(chan struct{}, workersNum)
	for {
		sem <- struct{}{}

		select {
		case <-ctx.Done():
			log.Info("message subscriber stopped")
			return
		default:
		}

		go func() {
			defer func() {
				<-sem
			}()

			msg, err := iter.Next()
			if err != nil {
				log.Errorf("failed to get next message: %+v", err)
				return
			}

			log.Infof("Received subject: %s\n", msg.Subject())

			handler, ok := handlers[msg.Subject()]
			if ok {
				if err := handler(ctx, msg); err != nil {
					log.Errorf("failed to handle message: %+v", err)
					return
				}

				msg.Ack()
			} else {
				log.Infof("handler not found subject: %+v", msg.Subject())
			}
		}()
	}
}
