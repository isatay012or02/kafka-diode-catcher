package main

import (
	"github.com/isatay012or02/kafka-diode-catcher/internal/adapters"
	"github.com/isatay012or02/kafka-diode-catcher/internal/application"
	"log"
)

func main() {

	hashCalculator := adapters.NewSHA1HashCalculator()

	udpReceiver, err := adapters.NewUDPReceiver(9999)
	if err != nil {
		log.Fatal(err)
	}

	kafkaWriter := adapters.NewKafkaWriter([]string{"localhost:9092"}, "example-topic")

	catcherService := application.NewCatcherService(udpReceiver, kafkaWriter, hashCalculator)

	err = catcherService.ReceiveAndPublishMessages()
	if err != nil {
		log.Fatal(err)
	}
}
