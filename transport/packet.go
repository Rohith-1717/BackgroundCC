package transport

import (
	"encoding/binary"
	"errors"
)

const (
	PacketTypeData uint8 = 1
	PacketTypeAck  uint8 = 2
)

type Packet struct {
	Type      uint8
	Seq       uint64
	Ack       uint64
	Timestamp int64
	Payload   []byte
}

func (p *Packet) MarshalBinary() ([]byte, error) {
	if p.Type != PacketTypeData && p.Type != PacketTypeAck {
		return nil, errors.New("invalid packet type")
	}
	payloadLen := len(p.Payload)
	buf := make([]byte, 1+8+8+8+4+payloadLen)
	buf[0] = p.Type
	binary.BigEndian.PutUint64(buf[1:9], p.Seq)
	binary.BigEndian.PutUint64(buf[9:17], p.Ack)
	binary.BigEndian.PutUint64(buf[17:25], uint64(p.Timestamp))
	binary.BigEndian.PutUint32(buf[25:29], uint32(payloadLen))
	copy(buf[29:], p.Payload)

	return buf, nil
}

func (p *Packet) UnmarshalBinary(b []byte) error {
	if len(b) < 29 {
		return errors.New("packet too short")
	}
	p.Type = b[0]
	if p.Type != PacketTypeData && p.Type != PacketTypeAck {
		return errors.New("invalid packet type")
	}
	p.Seq = binary.BigEndian.Uint64(b[1:9])
	p.Ack = binary.BigEndian.Uint64(b[9:17])
	p.Timestamp = int64(binary.BigEndian.Uint64(b[17:25]))
	payloadLen := binary.BigEndian.Uint32(b[25:29])
	if int(payloadLen) != len(b[29:]) {
		return errors.New("payload length mismatch")
	}
	if payloadLen > 0 {
		p.Payload = make([]byte, payloadLen)
		copy(p.Payload, b[29:])
	} else {
		p.Payload = nil
	}
	return nil
}
