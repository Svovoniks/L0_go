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
		fmt.Println("You are an idiot")
	}

	return string(enc)

}

func OrderFromMessage(message []byte) (*types.Order, error) {
	var jsonMap map[string]any
	err := json.Unmarshal(message, &jsonMap)

	if err != nil {
		return nil, err
	}

	if !types.IsValidOrder(jsonMap) {
		return nil, errors.New("Not a valid order")
	}

	if _, ok := jsonMap[types.OrderIdJsonKey].(string); !ok {
		return nil, errors.New(fmt.Sprint("Expected order_uid to be string, but got:", reflect.TypeOf(jsonMap[types.OrderIdJsonKey])))
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
	writer := kafka.Writer{
		Addr:     kafka.TCP("localhost:29092"),
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	}

	return &writer
}
func ProcessMessage(msg []byte, ctx *types.LocalContext) {
	order, err := OrderFromMessage(msg)

	if err != nil {
		fmt.Println("Received invalid order")
		fmt.Println(msg)
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
				fmt.Println("Conumer was stopped")
				break
			}
			fmt.Printf("Failed to read message: %s\n", err)
			continue
		}

		ProcessMessage(msg.Value, ctx)
		fmt.Println(ctx.Cache)
		fmt.Println(ctx.Db.GetAll())
	}

	reader.Close()
	ctx.WaitGroup.Done()
}

func WaitForExit() {
	var inp string

	for {
		fmt.Scanln(&inp)
		if inp == "exit" {
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
				fmt.Println("Producer was stopped")
				break
			}
			fmt.Printf("Couldn't write message: %s\n", err)
		}

		time.Sleep(time.Duration(time.Duration.Seconds(10)))
	}
	writer.Close()
	ctx.WaitGroup.Done()
}

func main() {
	readerCtx, cancelRReader := context.WithCancel(context.Background())
	writerCtx, cancelWReader := context.WithCancel(context.Background())
	ctx := types.GetLocalContext(&readerCtx, &writerCtx)

	ctx.WaitGroup.Add(1)
	go RunConsumerPipeline(&ctx)

	ctx.WaitGroup.Add(1)
	go RunProducerPipeline(&ctx)

	go ui.StartUI(&ctx)

	WaitForExit()

	cancelWReader()
	cancelRReader()

	ctx.WaitGroup.Wait()
}
