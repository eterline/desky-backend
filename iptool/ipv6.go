package iptool

import (
	"strconv"
	"strings"
)

var (
	ZeroIPv6      = NewIPv6(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0) // Default null address ip: ::/0
	LoopBackIPv6  = NewIPv6(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0) // Loopback/Localhost ip address: ::1
	BroadcastIPv6 = NewIPv6(255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 128)
)

func NewIPv6(
	a0, b0, c0, d0,
	a1, b1, c1, d1,
	a2, b2, c2, d2,
	a3, b3, c3, d3 byte,
	sub byte,
) *IP {
	return &IP{
		addr: []byte{
			a0, b0, c0, d0,
			a1, b1, c1, d1,
			a2, b2, c2, d2,
			a3, b3, c3, d3,
		},
		mask: sub,
		ver:  IPv6,
	}
}

// ParseIPv6 - parsing ip v6 addr string.
// Valid values example: 'c0ff:ffff:0a0a:0a0a:ffff:ffff:ff04::' or with subnet 'c0ff:ffff:0a0a:0a0a:ffff:ffff:ff04::/24'
func ParseIPv6(ip string) (addr *IP, ok bool) {
	addr = ZeroIPv6

	parts := strings.Split(ip, "/")

	if len(parts) > 2 {
		return nil, false
	}

	if parts[0] == "::" {
		return ZeroIPv6, true
	}

	if parts[0] == "::1" {
		return LoopBackIP, true
	}

	// replacing zero combinates
	ip = strings.ReplaceAll(parts[0], "::", ":0000")

	if len(parts) >= 1 {
		hexPairs := strings.Split(parts[0], ":")

		// b.c. because the length of slices in golang is a multiple of 3, with 8 elements = 9
		if len(hexPairs) != 9 {
			return nil, false
		}

		for idx := range 8 {

			// making bytes pairs
			pair, err := hexStringToDecimal(hexPairs[idx], 2)
			if err != nil {
				return nil, false
			}

			addr.addr[idx*2], addr.addr[idx*2+1] = pair[0], pair[1]
		}
	}

	if len(parts) > 1 {
		v, _ := strconv.Atoi(parts[1])
		if v > 128 {
			return nil, false
		}
		addr.mask = byte(v)
	}

	return addr, true
}
