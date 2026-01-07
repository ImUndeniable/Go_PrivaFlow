package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type EventProducer struct {
	writer *kafka.Writer
}

// NewEventProducer creates the connection to Kafka
// brokerUrl = "localhost:9092"
// topic = "erasure-requests"
func NewEventProducer(brokerUrl string, topic string) *EventProducer {
	// 👇 NEW: Try to create the topic automatically
	conn, err := kafka.Dial("tcp", brokerUrl)
	if err != nil {
		log.Println("⚠️  Failed to connect to Kafka for topic creation:", err)
	} else {
		defer conn.Close()

		err = conn.CreateTopics(kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})

		// We ignore "Topic already exists" errors, but log others
		if err != nil {
			log.Println("ℹ️  Topic setup info:", err)
		} else {
			log.Println("✅  Topic '" + topic + "' created successfully!")
		}
	}
	// 👆 END NEW CODE

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerUrl),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &EventProducer{writer: writer}
}

// Publish takes any data (value), turns it into JSON, and sends it
func (p *EventProducer) Publish(ctx context.Context, key string, value interface{}) error {
	// 1. Marshal struct to JSON
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// 2. Create the message
	msg := kafka.Message{
		Key:   []byte(key),
		Value: payload,
		Time:  time.Now(),
	}

	// 3. Send to Kafka
	return p.writer.WriteMessages(ctx, msg)
}

// Close disconnects
func (p *EventProducer) Close() {
	if err := p.writer.Close(); err != nil {
		log.Println("❌ Error closing Kafka writer:", err)
	}
}
