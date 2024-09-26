package test_utils_kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"net"
)


func GetKafkaClient() (*kafka.Client, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:29092")

	if err != nil {
		return nil, err
	}

	return &kafka.Client{
		Addr: addr,
	}, nil

}

func CreateTopic(name string, client *kafka.Client) error {
	topicCfg := kafka.TopicConfig{
		Topic:             name,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	rq := kafka.CreateTopicsRequest{
		Topics: []kafka.TopicConfig{topicCfg},
	}

	_, err := client.CreateTopics(context.Background(), &rq)

	return err
}

func DeleteTopic(name string, client *kafka.Client) error {
	deleteRq := kafka.DeleteTopicsRequest{
		Topics: []string{name},
	}

	_, err := client.DeleteTopics(context.Background(), &deleteRq)
	return err
}

func GetTestWriter() *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:29092"},
		Topic:    "test",
		Balancer: &kafka.LeastBytes{},
	})
}

func GetTestReader() *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:29092"},
		Topic:     "test",
		Partition: 0,
		MaxBytes:  10e6,
	})

	return reader
}
