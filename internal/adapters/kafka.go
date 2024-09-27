package adapters

import (
	"context"
	"github.com/isatay012or02/kafka-diode-catcher/internal/domain"
	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	writer *kafka.Writer
}

func NewKafkaWriter(brokers []string) *KafkaWriter {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
	})
	return &KafkaWriter{writer: writer}
}

func (kw *KafkaWriter) WriteMessage(message domain.Message) error {
	return kw.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: message.Topic,
		Key:   []byte(message.Key),
		Value: []byte(message.Value),
	})
}
