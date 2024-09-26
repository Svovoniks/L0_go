package kafka_reader

import (
	"context"
	"errors"
	"fmt"
	"l0/types/logger"

	"github.com/segmentio/kafka-go"
)

type KafkaReader struct {
	ctx     context.Context
	reader  *kafka.Reader
	Cancel  context.CancelFunc
	Stopped bool
}

type KafkaReaderContext struct {
	Host      string
	Port      string
	Topic     string
	Partition int
}

func GetKafkaReader(ctx KafkaReaderContext) *KafkaReader {
	readerCtx, cancelReader := context.WithCancel(context.Background())

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{fmt.Sprintf("%s:%s", ctx.Host, ctx.Port)},
		Topic:     ctx.Topic,
		Partition: ctx.Partition,
		MaxBytes:  10e6,
	})

	return &KafkaReader{
		ctx:     readerCtx,
		reader:  reader,
		Cancel:  cancelReader,
		Stopped: false,
	}
}

func (k *KafkaReader) ConsumeMessage() (*kafka.Message, error) {
	msg, err := k.reader.ReadMessage(k.ctx)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			logger.Logger.Info().
				Msg("Consumer was stopped")

			k.Stopped = true
		} else {
			logger.Logger.Warn().
				Msg("Error while retreiving message")

		}

		return nil, err
	}

	return &msg, nil

}

func (k *KafkaReader) Close() {
	k.reader.Close()
}
