package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

type Consumer struct {
	reader *kafka.Reader
}

type Message struct {
	Topic string
	Key   string
	Value interface{}
}

func NewProducer(brokers []string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{writer: writer}
}

func (p *Producer) Publish(ctx context.Context, topic, key string, value interface{}) error {
	message := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
	}

	if value != nil {
		valBytes, err := json.Marshal(value)
		if err != nil {
			return err
		}
		message.Value = valBytes
	}

	return p.writer.WriteMessages(ctx, message)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{reader: reader}
}

func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func CreateTopics(brokers []string, topics ...string) error {
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerAddr := fmt.Sprintf("%s:%d", controller.Host, controller.Port)
	connController, err := kafka.Dial("tcp", controllerAddr)
	if err != nil {
		return err
	}
	defer connController.Close()

	for _, topic := range topics {
		err = connController.CreateTopics(
			kafka.TopicConfig{
				Topic:             topic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
		)
		if err != nil {
			log.Printf("Failed to create topic %s: %v", topic, err)
		}
	}

	return nil
}
