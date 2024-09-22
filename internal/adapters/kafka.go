package adapters

import (
	"github.com/isatay012or02/kafka-diode-catcher/internal/domain"
	"github.com/segmentio/kafka-go"
)

type KafkaReader struct {
	reader *kafka.Reader
}

func NewKafkaReader(brokers []string, topic string, groupID string) *KafkaReader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &KafkaReader{reader: reader}
}

func (kr *KafkaReader) ReadMessage() (domain.Message, error) {
	msg, err := kr.reader.ReadMessage(nil)
	if err != nil {
		return domain.Message{}, err
	}

	return domain.Message{
		Data: string(msg.Value),
	}, nil
}

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
		Value: []byte(message.Data),
	})
}
