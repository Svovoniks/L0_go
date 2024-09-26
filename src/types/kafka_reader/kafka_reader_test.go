package kafka_reader

import (
	"context"
	"fmt"
	kafkaU "l0/test_utils/kafka"
	"l0/types/random"
	"testing"

	"github.com/segmentio/kafka-go"
)

func TestConsumeMessage(t *testing.T) {
	writer := kafkaU.GetTestWriter()
	defer writer.Close()

	client, err := kafkaU.GetKafkaClient()
	if err != nil {
		t.Error("Couldn't create kafka client")
		return
	}

	err = kafkaU.CreateTopic("test", client)
	if err != nil {
		t.Error("Couldn't create a topic")
		return
	}
	defer kafkaU.DeleteTopic("test", client)

	msg := random.RandomString(25)

	err = writer.WriteMessages(context.Background(), kafka.Message{Value: []byte(msg)})

	if err != nil {
		fmt.Println(err)
		t.Error("Couldn't produce a message")
	}

	reader := GetKafkaReader(KafkaReaderContext{
		Host:      "localhost",
		Port:      "29092",
		Partition: 0,
		Topic:     "test",
	})


	msgR, errC := reader.ConsumeMessage()
	if errC != nil {
		t.Error("Couldn't consume a message")
	}

	if string(msgR.Value) != msg {
		t.Error("Received the wrong message")
	}
}
