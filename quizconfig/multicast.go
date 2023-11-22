package quizconfig

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"

	"github.com/abaskin/testparts"
	"github.com/pkg/errors"
)

type TransportMulticast_t struct {
	sendConn *net.UDPConn
	recvConn *net.UDPConn
}

func NewMultiCast(recvChan chan<- *ClientData_t) (*TransportMulticast_t, error) {
	addr, err := net.ResolveUDPAddr("udp4", MultiCastAddress)
	if err != nil {
		return &TransportMulticast_t{}, errors.Wrap(err, "ResolveUDPAddr failed")
	}

	// Open up a connection
	rconn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return &TransportMulticast_t{}, errors.Wrap(err, "ListenMulticastUDP failed")
	}
	rconn.SetReadBuffer(MaxDatagramSize)

	sconn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return &TransportMulticast_t{}, errors.Wrap(err, "DialUDP failed")
	}

	go func(recvChan chan<- *ClientData_t, rconn *net.UDPConn) {
		for {
			buffer := make([]byte, MaxDatagramSize)
			msgLen, _, err := rconn.ReadFromUDP(buffer)
			switch {
			case err == nil:
				msg := &ClientData_t{}
				if err := gob.NewDecoder(bytes.NewBuffer(buffer[:msgLen])).
					Decode(msg); err != nil {
					log.Println("gob decode failed:", err)
					continue
				}
				recvChan <- msg
			case errors.Is(err, net.ErrClosed):
				return
			default:
				log.Println("ReadFromUDP failed:", err)
			}
		}
	}(recvChan, rconn)

	transport := &TransportMulticast_t{
		recvConn: rconn,
		sendConn: sconn,
	}
	transport.SendNoop()
	return transport, nil
}

func (m *TransportMulticast_t) Send(data *ClientData_t, topic string) error {
	log.Println(data.Action.String(), "->")
	var msg bytes.Buffer
	if err := gob.NewEncoder(&msg).Encode(data); err != nil {
		return errors.Wrap(err, "gob encode error")
	}

	if _, err := m.sendConn.Write(msg.Bytes()); err != nil {
		return errors.Wrap(err, "send failed")
	}
	return nil
}

func (m *TransportMulticast_t) SendNoop() error {
	return m.Send(&ClientData_t{
		Action:   ActionNoop,
		Question: &testparts.GormQuestion{},
	}, "")
}

func (m *TransportMulticast_t) ShutDown() {
	m.sendConn.Close()
	m.recvConn.Close()
}
