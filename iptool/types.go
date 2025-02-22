package iptool

type IPSocket struct {
	Port uint16
	IP   *IP
}

type IP struct {
	addr []byte
	mask byte
	ver  IPv
}

// IP address version v4/v6
type IPv byte

const (
	_ IPv = iota
	IPv4
	IPv6
)

func (i IPv) String() string {
	switch i {
	case IPv4:
		return "IPv4"
	case IPv6:
		return "IPv6"
	default:
		return "<nil>"
	}
}
