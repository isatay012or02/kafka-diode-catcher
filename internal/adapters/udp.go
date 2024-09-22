package adapters

import (
	"encoding/json"
	"github.com/isatay012or02/kafka-diode-catcher/internal/domain"
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

func (r *UDPReceiver) Receive() (domain.Message, error) {
	buf := make([]byte, 1024)
	n, _, err := r.conn.ReadFromUDP(buf)
	if err != nil {
		return domain.Message{}, err
	}

	var message domain.Message
	err = json.Unmarshal(buf[:n], &message)
	return message, err
}
