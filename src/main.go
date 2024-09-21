package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"l0/types"
	"l0/ui"
	"math/rand"
	"reflect"
	"runtime"
	"time"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}
	return string(result)
}

func RandomValidOrder() string {
	orderMap := make(map[string]string)

	for _, field := range types.RequiredOrderFields {
		orderMap[field] = generateRandomString(10)
	}

	enc, err := json.Marshal(orderMap)

	if err != nil {
		types.Logger.Fatal().Err(err).Msg("Shoud never happen")
	}

	return string(enc)

}

func OrderFromMessage(message []byte) (*types.Order, error) {
	var jsonMap map[string]any
	err := json.Unmarshal(message, &jsonMap)

	if err != nil {
		types.Logger.Warn().
			Str("kafka_message", string(message)).
			Msg("Received invalid Json")
		return nil, err
	}

	if !types.IsValidOrder(jsonMap) {
		types.Logger.Warn().
			Str("kafka_message", string(message)).
			Msg("Received invalid Order")
		return nil, errors.New("Not a valid order")
	}

	if _, ok := jsonMap[types.OrderIdJsonKey].(string); !ok {
		uid_err := errors.New(fmt.Sprint("Expected order_uid to be string, but got:", reflect.TypeOf(jsonMap[types.OrderIdJsonKey])))

		types.Logger.Warn().
			Str("kafka_message", string(message)).
			Err(uid_err).
			Msg("Received order with invalid 'order_uid'")
		return nil, uid_err
	}

	return &types.Order{
		Id:      jsonMap[types.OrderIdJsonKey].(string),
		JsonStr: string(message),
	}, nil

}

func ProcessOrder(order *types.Order, db *types.DB, cache *types.Cache) {
	db.Put(order)
	cache.Put(order)
}

func GetKafkaReader() *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:29092"},
		Topic:     "orders",
		Partition: 0,
		MaxBytes:  10e6,
	})

	return reader
}

func GetKafkaWriter() *kafka.Writer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:29092"},
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	})

	return writer
}
func ProcessMessage(msg []byte, ctx *types.LocalContext) {
	order, err := OrderFromMessage(msg)

	if err != nil {
		types.Logger.Warn().
			Err(err).
			Str("kafka_message", string(msg)).
			Msg("Skipping message")
		return
	}

	ctx.Db.Put(order)
	ctx.Cache.Put(order)
}

func RunConsumerPipeline(ctx *types.LocalContext) {
	reader := GetKafkaReader()

	for {
		msg, err := reader.ReadMessage(*ctx.Reader_ctx)

		if err != nil {
			if errors.Is(err, context.Canceled) {
				types.Logger.Info().
					Msg("Consumer was stopped")
				fmt.Println("Consumer was stopped")
				break
			}

			types.Logger.Info().
				Err(err).
				Msg("Couldn't read message, skipping")
			continue
		}

		ProcessMessage(msg.Value, ctx)
	}

	reader.Close()
	ctx.WaitGroup.Done()
}

func WaitForExit() {
	var inp string

	for {
		fmt.Scanln(&inp)
		if inp == "exit" {
			types.Logger.Info().
				Msg("Received exit request")
			return
		}
	}
}

func RunProducerPipeline(ctx *types.LocalContext) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:29092"),
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	}

	for {
		err := writer.WriteMessages(*ctx.Writer_ctx,
			kafka.Message{
				Value: []byte(RandomValidOrder()),
			},
		)

		if err != nil {
			if errors.Is(err, context.Canceled) {
				types.Logger.Info().
					Msg("Producer was stopped")
				fmt.Println("Producer was stopped")
				break
			}
			types.Logger.Warn().
				Err(err).
				Msg("Couldn't write message")
		}

		time.Sleep(time.Duration(time.Duration.Seconds(5)))
	}
	writer.Close()
	ctx.WaitGroup.Done()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if writer, err := types.SetupLogger(); err == nil {
		defer writer.Close()
	}

	readerCtx, cancelReader := context.WithCancel(context.Background())
	writerCtx, cancelWriter := context.WithCancel(context.Background())
	ctx := types.GetLocalContext(&readerCtx, &writerCtx)

	ctx.WaitGroup.Add(1)
	go RunConsumerPipeline(&ctx)

	ctx.WaitGroup.Add(1)
	go RunProducerPipeline(&ctx)

	go ui.StartUI(&ctx)

	WaitForExit()

	cancelReader()
	cancelWriter()

	ctx.WaitGroup.Wait()
	ctx.Db.Db.Close()
}
