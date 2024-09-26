package main

import (
	"fmt"
	"l0/types/config"
	"l0/types/kafka_reader"
	"l0/types/kafka_writer"
	local_context "l0/types/local_context"
	"l0/types/logger"
	_order "l0/types/order"
	"l0/ui"
	"runtime"
	"time"

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
		}

		ProcessMessage(msg.Value, ctx)
	}

	logger.Logger.Info().
		Msg("Shutting down consumer pipeline")

	ctx.Reader.Close()
	ctx.WaitGroup.Done()
}

func RunProducerPipeline(ctx *local_context.LocalContext) {
	for {
		err := ctx.Writer.ProduceMessage(_order.RandomValidOrder())

		if err != nil {
			if ctx.Writer.Stopped {
				break
			}
			logger.Logger.Warn().
				Err(err).
				Msg("Skipping message")
		}

		time.Sleep(time.Duration(time.Duration.Seconds(10)))
	}

	logger.Logger.Info().
		Msg("Shutting down producer pipeline")
	fmt.Println("Shutting down producer pipeline")

	ctx.Writer.Close()
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

	writer := kafka_writer.GetKafkaWriter(kafka_writer.KafkaWriterContext{
		Host: config.KafkaHost,
		Port: config.KafkaPort,
        Topic: config.KafkaTopic,
	})

	ctx, errC := local_context.GetLocalContext(config, reader, writer)
	if errC != nil {
		logger.Logger.Warn().
			Msg("Couldn't get local context, exiting...")
		return
	}

	ctx.WaitGroup.Add(1)
	go RunConsumerPipeline(ctx)

	ctx.WaitGroup.Add(1)
	go RunProducerPipeline(ctx)

	go ui.StartUI(ctx)

	WaitForExit()

	writer.Cancel()
	reader.Cancel()

	ctx.WaitGroup.Wait()
	ctx.Db.Db.Close()
}
