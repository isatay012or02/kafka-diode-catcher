package main

import (
	"fmt"
	"github.com/isatay012or02/kafka-diode-catcher/config"
	"github.com/isatay012or02/kafka-diode-catcher/internal/adapters"
	"github.com/isatay012or02/kafka-diode-catcher/internal/application"
	"github.com/isatay012or02/kafka-diode-catcher/internal/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg, err := config.Init("config.json")
	if err != nil {
		panic(err)
	}

	go func(cfg *config.Config) {
		hashCalculator := adapters.NewSHA1HashCalculator()
		kafkaWriter := adapters.NewKafkaWriter(cfg.Queue.Brokers)
		udpReceiver, err := adapters.NewUDPReceiver(cfg.UdpAddress.Ip, cfg.UdpAddress.Port)
		if err != nil {
			panic(err)
		}

		catcherService := application.NewCatcherService(udpReceiver, kafkaWriter, hashCalculator, cfg.Queue)

		err = catcherService.ReceiveAndPublishMessages()
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
