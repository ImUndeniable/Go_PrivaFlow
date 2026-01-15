package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type EventConsumer struct {
	reader *kafka.Reader
}

func NewEventConsumer(brokerUrl string, topic string, groupID string) *EventConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerUrl},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return &EventConsumer{reader: reader}
}

func (c *EventConsumer) FetchMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *EventConsumer) Close() {
	if err := c.reader.Close(); err != nil {
		log.Println("❌ Error closing Kafka reader:", err)
	}
}
