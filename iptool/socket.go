package iptool

import "fmt"

func NewSocket(ip *IP, port int) (a *IPSocket, ok bool) {
	if port < 0 || port > 65535 {
		return nil, false
	}

	return &IPSocket{
		Port: uint16(port),
		IP:   ip,
	}, true
}

func (so *IPSocket) String() string {
	switch so.IP.ver {
	case IPv4:
		return fmt.Sprintf("%s:%v", so.IP, so.Port)
	case IPv6:
		return fmt.Sprintf("[%s]:%v", so.IP, so.Port)
	default:
		panic(ErrUncompVersion)
	}
}
