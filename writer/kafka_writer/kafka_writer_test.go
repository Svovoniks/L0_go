package kafka_writer

import (
	"context"
	kafkaU "writer/test_utils/kafka"
	"writer/random"
	"testing"
)

func TestProduceMessage(t *testing.T) {
	reader := kafkaU.GetTestReader()
	defer reader.Close()

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

	writer := GetKafkaWriter(KafkaWriterContext{
		Host: "localhost",
		Port: "29092",
        Topic: "test",
	})

	err = writer.ProduceMessage(msg)

	if err != nil {
		t.Error("Couldn't produce a message")
	}

	msgR, errC := reader.ReadMessage(context.Background())
	if errC != nil {
		t.Error("Couldn't read the message")
	}

	if string(msgR.Value) != msg {
		t.Error("Received the wrong message")
	}
}
