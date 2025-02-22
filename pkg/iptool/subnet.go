package iptool

import (
	"math"
)

// Subnet value number
func (i *IP) SubnetValue() int {
	return int(i.mask)
}

func (i *IP) Subnet() [4]byte {
	m := i.Mask()
	v := [4]byte{}

	for idx, val := range m {
		if val != 0 {
			v[idx] = i.addr[idx]
		}
	}

	return v
}

// IncMask - increments mask value: addr/23 -> addr/24
func (i *IP) IncMask() {
	if ok := i.safeSubnetSet(i.mask + 1); !ok {
		panic("uncorrect subnet value")
	}
}

// DecMask - decrements mask value: addr/23 -> addr/22
func (i *IP) DecMask() {
	if ok := i.safeSubnetSet(i.mask - 1); !ok {
		panic("uncorrect subnet value")
	}
}

func (i *IP) MustMaskValue(value int) {
	if ok := i.safeSubnetSet(byte(value)); !ok {
		panic("uncorrect subnet value")
	}
}

func (i *IP) Mask() []byte {
	if i.ver != (IPv4) {
		panic(ErrUncompVersion)
	}

	mask := ZeroIP.addr

	m := i.mask

	switch i.mask {

	case 0:
		break
	case 32:
		return BroadcastIP.addr

	case 128:
		return BroadcastIPv6.addr

	default:
		for i := range mask {
			rt := 0
			for rt < 8 {
				if m == 0 {
					return mask
				}
				m--
				rt++
				mask[i] = mask[i] + byte(math.Pow(2, float64(8-rt)))
			}
		}
	}

	return mask
}
