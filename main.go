package main

import (
	"fmt"
	"github.com/isatay012or02/kafka-diode-catcher/config"
	"github.com/isatay012or02/kafka-diode-catcher/internal/adapters"
	"github.com/isatay012or02/kafka-diode-catcher/internal/application"
	"github.com/isatay012or02/kafka-diode-catcher/internal/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {

	cfg, err := config.Init("config.json")
	if err != nil {
		panic(err)
	}

	go func(cfg *config.Config) {
		var loggerBroker []string
		loggerBroker[0] = os.Getenv("KAFKA_LOGGER_BROKER")
		loggerTopic := os.Getenv("KAFKA_LOGGER_TOPIC")

		logger := adapters.NewKafkaLogger(loggerBroker, loggerTopic)

		enableHashEnv := os.Getenv("ENABLE_HASH")
		enableHash := false
		if enableHashEnv == "true" {
			enableHash = true
		}

		udpIP := os.Getenv("UPD_IP")
		updPortEnv := os.Getenv("UPD_PORT")
		udpPort, err := strconv.Atoi(updPortEnv)
		if err != nil {
			panic(err)
		}

		hashCalculator := adapters.NewSHA1HashCalculator()
		kafkaWriter := adapters.NewKafkaWriter(cfg.Queue.Brokers)
		udpReceiver, err := adapters.NewUDPReceiver(udpIP, udpPort)
		if err != nil {
			panic(err)
		}

		catcherService := application.NewCatcherService(udpReceiver, kafkaWriter, hashCalculator, cfg.Queue, enableHash, logger)

		err = catcherService.ReceiveAndPublishMessages()
		logger.SendMetricsToKafka()
		if err != nil {
			panic(err)
		}
	}(cfg)

	srv, err := http.NewServer(cfg)
	if err != nil {
		panic(err)
	}

	startServerErrorCH := srv.Start()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err = <-startServerErrorCH:
		{
			panic(err)
		}
	case q := <-quit:
		{
			fmt.Printf("receive signal %s, stopping server...\n", q.String())
			if err = srv.Stop(); err != nil {
				fmt.Printf("stop server error: %s\n", err.Error())
			}
		}
	}
}
