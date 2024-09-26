package main

import (
	"fmt"
	"l0/types/config"
	"l0/types/kafka_reader"
	local_context "l0/types/local_context"
	"l0/types/logger"
	_order "l0/types/order"
	"l0/ui"
	"runtime"

	_ "github.com/lib/pq"
)

func ProcessMessage(msg []byte, ctx *local_context.LocalContext) {
	order, err := _order.OrderFromMessage(msg)

	if err != nil {
		logger.Logger.Warn().
			Err(err).
			Str("kafka_message", string(msg)).
			Msg("Skipping message")
		return
	}

	ctx.Db.Put(order)
	ctx.Cache.Put(order)
}
func RunConsumerPipeline(ctx *local_context.LocalContext) {
	for {
		msg, err := ctx.Reader.ConsumeMessage()

		if err != nil {
			if ctx.Reader.Stopped {
				break
			}
			logger.Logger.Warn().
				Err(err).
				Msg("Skipping message")
            continue
		}

		ProcessMessage(msg.Value, ctx)
	}

	logger.Logger.Info().
		Msg("Shutting down consumer pipeline")

	ctx.Reader.Close()
	ctx.WaitGroup.Done()
}

func WaitForExit() {
	var inp string

	for {
		fmt.Scanln(&inp)
		if inp == "exit" {
			logger.Logger.Info().
				Msg("Received exit request")
			return
		}
	}
}

func main() {
	config, err := config.GetConfig()
	if err != nil {
		fmt.Println("No config, exiting...")
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	if writer, err := logger.SetupLogger(); err == nil {
		defer writer.Close()
	}

	reader := kafka_reader.GetKafkaReader(kafka_reader.KafkaReaderContext{
		Host:      config.KafkaHost,
		Port:      config.KafkaPort,
		Topic:     config.KafkaTopic,
		Partition: 0,
	})

	ctx, errC := local_context.GetLocalContext(config, reader)
	if errC != nil {
		logger.Logger.Warn().
			Msg("Couldn't get local context, exiting...")
		return
	}

	ctx.WaitGroup.Add(1)
	go RunConsumerPipeline(ctx)

	go ui.StartUI(ctx)

	WaitForExit()

	reader.Cancel()

	ctx.WaitGroup.Wait()
	ctx.Db.Db.Close()
}
