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
	"sync"
	"syscall"
	"time"
)

func main() {

	cfg, err := config.Init("config.json")
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup

	wg.Add(1)

	go func(cfg *config.Config) {
		defer wg.Done()

		loggerTopic := os.Getenv("KAFKA_LOGGER_TOPIC")

		enableHashEnv := os.Getenv("ENABLE_HASH")
		enableHash := false
		if enableHashEnv == "true" {
			enableHash = true
		}

		udpIP := os.Getenv("UDP_IP")
		updPortEnv := os.Getenv("UDP_PORT")
		udpPort, err := strconv.Atoi(updPortEnv)
		if err != nil {
			panic(err)
		}

		hashCalculator := adapters.NewSHA1HashCalculator()
		kafkaWriter := adapters.NewKafkaWriter(cfg.Queue.Brokers, loggerTopic)

		kafkaWriter.Log(fmt.Sprintf("[%v][INFO]Caster service started", time.Now()))

		udpReceiver, err := adapters.NewUDPReceiver(udpIP, udpPort)
		if err != nil {
			panic(err)
		}

		catcherService := application.NewCatcherService(udpReceiver, kafkaWriter, hashCalculator, cfg.Queue, enableHash)
		err = catcherService.ReceiveAndPublishMessages()
		kafkaWriter.SendMetricsToKafka()
		kafkaWriter.Close()
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

	wg.Wait()
}
