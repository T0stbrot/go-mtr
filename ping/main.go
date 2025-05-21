package ping

import (
	"fmt"
	"net"
	"math/rand"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type PingResult struct {
	Target  string `json:"target"`
	LastHop string `json:"lasthop"`
	RTT     string `json:"rtt,omitempty"`
	Message string `json:"message,omitempty"`
}

type PingProto struct {
	IP string
	IPF string
	Listen string
	Type icmp.Type
	MessageType int
	Conn4 *ipv4.PacketConn
	Conn6 *ipv6.PacketConn
}

func Ping(ver int, destination string, ttl int, timeout int, seq int) PingResult {
	result := PingResult{Target: destination}

	proto := PingProto{IP: "ip4", IPF: "ip4:icmp", Listen: "0.0.0.0", Type: ipv4.ICMPTypeEcho, MessageType: 1}

	if ver == 6 {
		proto.IP = "ip6"
		proto.IPF = "ip6:ipv6-icmp"
		proto.Listen = "::"
		proto.Type = ipv6.ICMPTypeEchoRequest
		proto.MessageType = 58
	}
	
	conn, err := net.ListenPacket(fmt.Sprintf("%s", proto.IPF), proto.Listen)
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}
	defer conn.Close()

	if ver == 6 {
		proto.Conn6 = ipv6.NewPacketConn(conn)
		proto.Conn6.SetHopLimit(ttl)
	} else {
		proto.Conn4 = ipv4.NewPacketConn(conn)
		proto.Conn4.SetTTL(ttl)
	}

	dst, err := net.ResolveIPAddr(proto.IP, destination)
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	id := int(rand.Intn(65536))

	icmpMessage := icmp.Message{
		Type: proto.Type,
		Code: 0,
		Body: &icmp.Echo{
			ID:   id,
			Seq:  seq,
			Data: make([]byte, 16),
		},
	}

	msgBytes, err := icmpMessage.Marshal(nil)
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	sT := time.Now()

	if _, err := conn.WriteTo(msgBytes, dst); err != nil {
		result.Message = fmt.Sprintf("%v", err)
		return result
	}

	buf := make([]byte, 1280)
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))

	n, addr, err := conn.ReadFrom(buf)
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
	}

	eT := time.Now()
	result.RTT = fmt.Sprintf("%.3f", float64(eT.Sub(sT).Microseconds())/1000)

	reply, err := icmp.ParseMessage(proto.MessageType, buf[:n])
	if err != nil {
		result.Message = fmt.Sprintf("%v", err)
	}

	result.LastHop = addr.String()
	switch reply.Type {
	case ipv4.ICMPTypeEchoReply:
		result.Message = "suceed"
	case ipv4.ICMPTypeTimeExceeded:
		result.Message = "timeexceed"
	case ipv6.ICMPTypeEchoReply:
		result.Message = "succeed"
	case ipv6.ICMPTypeTimeExceeded:
		result.Message = "timeexceed"
	default:
		result.Message = fmt.Sprintf("%v", reply)
	}

	return result
}
