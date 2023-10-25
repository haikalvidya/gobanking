package nats

import (
	"context"
	"gobanking/pkg/logger"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func NewJetStream(ctx context.Context, nc *nats.Conn, streamName string, subjects []string, log logger.Logger) (*jetstream.JetStream, *jetstream.Stream, error) {
	js, err := jetstream.New(nc)
	if err != nil {
		log.Errorf("failed to create jetstream: %v", err)
		return nil, nil, err
	}

	stream, _ := js.Stream(ctx, streamName)
	if stream == nil {
		stream, err = js.CreateStream(ctx, jetstream.StreamConfig{
			Name:      streamName,
			Subjects:  subjects,
			Retention: jetstream.WorkQueuePolicy,
			MaxAge:    7 * 24 * time.Hour,   // max age of stored messages is 7 days
			Discard:   jetstream.DiscardOld, // when the stream is full, discard old messages
		})
		if err != nil {
			log.Errorf("failed to create stream: %v", err)
			return nil, nil, err
		}
	} else {
		stream, err = js.UpdateStream(ctx, jetstream.StreamConfig{
			Name:      streamName,
			Subjects:  subjects,
			Retention: jetstream.WorkQueuePolicy,
			MaxAge:    7 * 24 * time.Hour,   // max age of stored messages is 7 days
			Discard:   jetstream.DiscardOld, // when the stream is full, discard old messages
		})
		if err != nil {
			log.Errorf("failed to update stream: %v", err)
			return nil, nil, err
		}
	}

	return &js, &stream, nil
}
