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
}

func NewCatcherService(udpReceiver *adapters.UDPReceiver, kafkaWriter *adapters.KafkaWriter,
	calculator ports.MessageHashCalculator, cfg config.Queue, enableHash bool) *CatcherService {

	return &CatcherService{
		UDPReceiver:    udpReceiver,
		KafkaWriter:    kafkaWriter,
		HashCalculator: calculator,
		Cfg:            cfg,
		EnableHash:     enableHash,
	}
}

func (c *CatcherService) ReceiveAndPublishMessages() error {

	timeStart := time.Now()

	for {
		msg, err := c.UDPReceiver.Receive()
		if err != nil {
			adapters.BroadcastStatus(-2, msg.Topic, "ERROR", time.Since(timeStart))
			c.KafkaWriter.Log(fmt.Sprintf("[%v][Error] %v", time.Now(), err.Error()))
			return err
		}

		if c.EnableHash {
			calculatedHash := c.HashCalculator.Calculate(msg.Value)
			if calculatedHash != msg.Hash {
				adapters.BroadcastStatusInc(-3, msg.Topic, "ERROR")
				c.KafkaWriter.Log(fmt.Sprintf("[%v][Error] %v, hash:%v, key: %v, value:%v", time.Now(), "hash mismatch", msg.Hash, msg.Key, msg.Value))
				return fmt.Errorf("hash mismatch")
			}
		}

		err = c.KafkaWriter.WriteMessage(msg)
		if err != nil {
			adapters.BroadcastStatus(-1, msg.Topic, "ERROR", time.Since(timeStart))
			c.KafkaWriter.Log(fmt.Sprintf("[%v][Error] %v", time.Now(), err.Error()))
			return err
		}

		adapters.BroadcastStatus(0, msg.Topic, "SUCCESS", time.Since(timeStart))
		c.KafkaWriter.SendMetricsToKafka()
	}
}
