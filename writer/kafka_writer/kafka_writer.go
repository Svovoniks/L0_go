package kafka_writer

import (
	"context"
	"errors"
	"fmt"
    "writer/logger"

	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	ctx     context.Context
	writer  *kafka.Writer
	Cancel  context.CancelFunc
	Stopped bool
}

type KafkaWriterContext struct {
	Host  string
	Port  string
	Topic string
}

func GetKafkaWriter(wctx KafkaWriterContext) *KafkaWriter {
	ctx, cancel := context.WithCancel(context.Background())

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{fmt.Sprintf("%s:%s", wctx.Host, wctx.Port)},
		Topic:    wctx.Topic,
		Balancer: &kafka.LeastBytes{},
	})

	return &KafkaWriter{
		ctx:     ctx,
		writer:  writer,
		Cancel:  cancel,
		Stopped: false,
	}
}

func (k *KafkaWriter) ProduceMessage(msg string) error {
	err := k.writer.WriteMessages(k.ctx,
		kafka.Message{
			Value: []byte(msg),
		},
	)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			logger.Logger.Info().
				Msg("Producer was stopped")
			k.Stopped = true
		} else {
			logger.Logger.Warn().
				Err(err).
				Msg("Couldn't write message")
		}

		return err
	}

	return nil
}

func (k *KafkaWriter) Close() {
	k.writer.Close()
}
