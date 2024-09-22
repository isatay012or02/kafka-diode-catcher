package events

import (
	"encoding/json"
	"github.com/isatay012or02/kafka-diode-catcher/models"
	"net"
)

type UDPReceiver struct {
	conn *net.UDPConn
}

func NewUDPReceiver(port int) (*UDPReceiver, error) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return nil, err
	}
	return &UDPReceiver{conn: conn}, nil
}

func (r *UDPReceiver) Receive() (models.Message, error) {
	buf := make([]byte, 1024)
	n, _, err := r.conn.ReadFromUDP(buf)
	if err != nil {
		return models.Message{}, err
	}

	var message models.Message
	err = json.Unmarshal(buf[:n], &message)
	return message, err
}
