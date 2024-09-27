package main

import (
	"fmt"
	"sync"
	"time"
	"writer/kafka_writer"
	"writer/logger"
	_order "writer/order"
)

func RunProducerPipeline(writer *kafka_writer.KafkaWriter, wg *sync.WaitGroup) {
	for {
        msg := _order.RandomValidOrder()
		err := writer.ProduceMessage(msg)

		if err != nil {
			if writer.Stopped {
				break
			}
			logger.Logger.Warn().
				Err(err).
				Msg("Skipping message")
		} else {
			logger.Logger.Info().
                Str("message", msg).
				Msg("Wrote message")
		}

		time.Sleep(5 * time.Second)
	}

	logger.Logger.Info().
		Msg("Shutting down producer pipeline")
	fmt.Println("Shutting down producer pipeline")

	writer.Close()
	wg.Done()
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
	logger.SetupLogger()
	writer := kafka_writer.GetKafkaWriter(kafka_writer.KafkaWriterContext{
		Host:  "localhost",
		Port:  "29092",
		Topic: "orders",
	})

	var wg sync.WaitGroup

	wg.Add(1)
	go RunProducerPipeline(writer, &wg)

	WaitForExit()

	writer.Cancel()

	wg.Wait()
}
