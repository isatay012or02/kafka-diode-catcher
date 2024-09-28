package adapters

import (
	"bytes"
	"context"
	"github.com/isatay012or02/kafka-diode-catcher/internal/domain"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaWriter struct {
	writer      *kafka.Writer
	loggerTopic string
}

func NewKafkaWriter(brokers []string, loggerTopic string) *KafkaWriter {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
	})
	return &KafkaWriter{writer: writer, loggerTopic: loggerTopic}
}

func (kw *KafkaWriter) WriteMessage(message domain.Message) error {
	return kw.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: message.Topic,
		Key:   []byte(message.Key),
		Value: []byte(message.Value),
	})
}

func (kl *KafkaWriter) Log(message string) {
	err := kl.WriteMessage(domain.Message{Key: "log-key", Topic: kl.loggerTopic, Value: message})

	if err != nil {
		panic(err)
	}
}

func (kl *KafkaWriter) SendMetricsToKafka() {
	// Сбор всех зарегистрированных метрик
	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		log.Println("Ошибка сбора метрик:", err)
		return
	}

	var buf bytes.Buffer
	encoder := expfmt.NewEncoder(&buf, expfmt.OpenMetricsType)

	for _, mf := range metricFamilies {
		err := encoder.Encode(mf)
		if err != nil {
			log.Println("Ошибка кодирования метрик:", err)
			return
		}
	}

	// Отправляем метрики в Kafka
	err = kl.WriteMessage(domain.Message{
		Key:   "metrics",
		Topic: kl.loggerTopic,
		Value: buf.String(),
	})

	if err != nil {
		panic(err)
	}
}

func (kl *KafkaWriter) Close() {
	kl.writer.Close()
}
