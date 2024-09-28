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
	Cfg            config.Queue
	EnableHash     bool
	Logger         *adapters.KafkaLogger
}

func NewCatcherService(udpReceiver *adapters.UDPReceiver, kafkaWriter *adapters.KafkaWriter,
	calculator ports.MessageHashCalculator, cfg config.Queue, enableHash bool,
	kafkaLogger *adapters.KafkaLogger) *CatcherService {

	return &CatcherService{
		UDPReceiver:    udpReceiver,
		KafkaWriter:    kafkaWriter,
		HashCalculator: calculator,
		Cfg:            cfg,
		EnableHash:     enableHash,
		Logger:         kafkaLogger,
	}
}

func (c *CatcherService) ReceiveAndPublishMessages() error {

	timeStart := time.Now()

	for {
		msg, err := c.UDPReceiver.Receive()
		if err != nil {
			adapters.BroadcastStatus(-2, msg.Topic, "ERROR", time.Since(timeStart))
			c.Logger.Log(fmt.Sprintf("[%v][Error] %v", time.Now(), err.Error()))
			return err
		}

		if c.EnableHash {
			calculatedHash := c.HashCalculator.Calculate(msg.Value)
			if calculatedHash != msg.Hash {
				adapters.BroadcastStatusInc(-3, msg.Topic, "ERROR")
				c.Logger.Log(fmt.Sprintf("[%v][Error] %v", time.Now(), "hash mismatch"))
				return fmt.Errorf("hash mismatch")
			}
		}

		err = c.KafkaWriter.WriteMessage(msg)
		if err != nil {
			adapters.BroadcastStatus(-1, msg.Topic, "ERROR", time.Since(timeStart))
			c.Logger.Log(fmt.Sprintf("[%v][Error] %v", time.Now(), err.Error()))
			return err
		}

		adapters.BroadcastStatus(0, msg.Topic, "SUCCESS", time.Since(timeStart))
		c.Logger.SendMetricsToKafka()
	}
}
