package transport

func (rt *ReliableTransport) ReceiveLoop() {
	buf := make([]byte, 2048)
	for {
		n, addr, err := rt.UDP.Receive(buf)
		if err != nil {
			return
		}
		var pkt Packet
		if pkt.UnmarshalBinary(buf[:n]) != nil {
			continue
		}
		rt.HandleIncoming(&pkt, addr)
	}
}
