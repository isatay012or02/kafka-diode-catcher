package application

import (
	"fmt"
	"github.com/isatay012or02/kafka-diode-catcher/internal/adapters"
	"github.com/isatay012or02/kafka-diode-catcher/internal/ports"
	"time"
)

type CatcherService struct {
	UDPReceiver    *adapters.UDPReceiver
	KafkaWriter    *adapters.KafkaWriter
	HashCalculator ports.MessageHashCalculator
}

func NewCatcherService(udpReceiver *adapters.UDPReceiver, kafkaWriter *adapters.KafkaWriter,
	calculator ports.MessageHashCalculator) *CatcherService {

	return &CatcherService{
		UDPReceiver:    udpReceiver,
		KafkaWriter:    kafkaWriter,
		HashCalculator: calculator,
	}
}

func (c *CatcherService) ReceiveAndPublishMessages() error {

	timeStart := time.Now()

	for {
		msg, err := c.UDPReceiver.Receive()
		if err != nil {
			return err
		}

		calculatedHash := c.HashCalculator.Calculate(msg.Data)
		if calculatedHash != msg.Hash {
			return fmt.Errorf("hash mismatch")
		}

		err = c.KafkaWriter.WriteMessage(msg)
		if err != nil {
			adapters.BroadcastStatus(-1, msg.Topic, "ERROR", time.Since(timeStart))
			return err
		}

		adapters.BroadcastStatus(0, msg.Topic, "SUCCESS", time.Since(timeStart))
	}
}
