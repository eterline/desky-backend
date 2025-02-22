package iptool

import (
	"strconv"
	"strings"
)

var (
	ZeroIP      = NewIP(0, 0, 0, 0, 0)         // Default null address ip: 0.0.0.0/0
	LoopBackIP  = NewIP(127, 0, 0, 1, 8)       // Loopback/Localhost ip address: 127.0.0.1
	LinkLocalIP = NewIP(169, 254, 0, 0, 16)    // Automatic IP configuration in the absence of DHCP: 169.254.0.0/16
	BroadcastIP = NewIP(255, 255, 255, 255, 0) // Broadcasting to all devices on the local network: 255.255.255.255
	APIPA       = NewIP(169, 254, 0, 0, 16)    // Automatic private IP addressing: 169.254.0.0/16
)

// RFC 1918
// Address Allocation for Private Internets
// February 1996
var (
	// PrivateClassA RFC 1918 - 10.xxx.xxx.xxx/8  (10/8 prefix)
	PrivateClassA = NewIP(10, 0, 0, 0, 8)
	// PrivateClassB RFC 1918 - 172.31.xxx.xxx/12  (172.16/12 prefix)
	PrivateClassB = NewIP(172, 31, 0, 0, 12)
	// PrivateClassC RFC 1918 - 192.168.xxx.xxx/16 (192.168/16 prefix)
	PrivateClassC = NewIP(192, 168, 0, 0, 16)
)

func NewIP(a0, b0, c0, d0 byte, sub byte) *IP {
	return &IP{
		addr: []byte{a0, b0, c0, d0},
		mask: sub,
		ver:  IPv4,
	}
}

// Parse ip addr string.
// Valid values example: '10.192.10.2' or with subnet '10.192.10.2/24'
func ParseIPv4(ip string) (addr *IP, ok bool) {
	addr = ZeroIP

	parts := strings.Split(ip, "/")

	if len(parts) > 2 {
		return nil, false
	}

	if len(parts) >= 1 {
		parties := strings.Split(parts[0], ".")

		if len(parties) != 4 {
			return nil, false
		}

		for idx, v := range parties {
			p, err := strconv.Atoi(v)

			if err != nil || p > 255 || p < 0 {
				return nil, false
			}
			addr.addr[idx] = byte(p)
		}
	}

	if len(parts) > 1 {
		v, _ := strconv.Atoi(parts[1])
		if v > 32 {
			return nil, false
		}
		addr.mask = byte(v)
	}

	return addr, true
}

// end
