package transport

import (
	"net"
)

type UDPTransport struct {
	conn *net.UDPConn
}

func NewUDPTransport(localAddr string) (*UDPTransport, error) {
	addr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	return &UDPTransport{
		conn: conn,
	}, nil
}

func (u *UDPTransport) Send(addr *net.UDPAddr, pkt *Packet) error {
	b, err := pkt.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = u.conn.WriteToUDP(b, addr)
	return err
}

func (u *UDPTransport) Receive(buf []byte) (int, *net.UDPAddr, error) {
	n, addr, err := u.conn.ReadFromUDP(buf)
	if err != nil {
		return 0, nil, err
	}
	return n, addr, nil
}

func (u *UDPTransport) Close() error {
	return u.conn.Close()
}
