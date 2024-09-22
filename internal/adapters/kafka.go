package adapters

import (
	"github.com/isatay012or02/kafka-diode-catcher/internal/domain"
	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	writer *kafka.Writer
}

func NewKafkaWriter(brokers []string, topic string) *KafkaWriter {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	return &KafkaWriter{writer: writer}
}

func (kw *KafkaWriter) WriteMessage(message domain.Message) error {
	return kw.writer.WriteMessages(nil, kafka.Message{
		Topic: message.Topic,
		Value: []byte(message.Data),
	})
}
