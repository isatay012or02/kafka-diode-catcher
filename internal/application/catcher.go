package application

import (
	"fmt"
	"github.com/isatay012or02/kafka-diode-catcher/config"
	"github.com/isatay012or02/kafka-diode-catcher/internal/adapters"
	"github.com/isatay012or02/kafka-diode-catcher/internal/ports"
	"time"
)

type CatcherService struct {
	UDPReceiver    *adapters.UDPReceiver
	KafkaWriter    *adapters.KafkaWriter
	HashCalculator ports.MessageHashCalculator
	cfg            config.Queue
}

func NewCatcherService(udpReceiver *adapters.UDPReceiver,
	calculator ports.MessageHashCalculator, cfg config.Queue) *CatcherService {

	return &CatcherService{
		UDPReceiver:    udpReceiver,
		HashCalculator: calculator,
		cfg:            cfg,
	}
}

func (c *CatcherService) ReceiveAndPublishMessages() error {

	timeStart := time.Now()

	for {
		msg, err := c.UDPReceiver.Receive()
		if err != nil {
			adapters.BroadcastStatus(-2, msg.Topic, "ERROR", time.Since(timeStart))
			adapters.BroadcastStatusInc(-2, msg.Topic, "ERROR")
			return err
		}

		calculatedHash := c.HashCalculator.Calculate(msg.Data)
		if calculatedHash != msg.Hash {
			return fmt.Errorf("hash mismatch")
		}

		kafkaWriter := adapters.NewKafkaWriter(c.cfg.Brokers, msg.Topic)

		err = kafkaWriter.WriteMessage(msg)
		if err != nil {
			adapters.BroadcastStatus(-1, msg.Topic, "ERROR", time.Since(timeStart))
			adapters.BroadcastStatusInc(-1, msg.Topic, "ERROR")
			return err
		}

		adapters.BroadcastStatus(0, msg.Topic, "SUCCESS", time.Since(timeStart))
	}
}
