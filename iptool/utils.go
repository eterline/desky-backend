package iptool

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// String - displays ip addr in string format.
// 10.192.10.1 or c0ff:ffff:0a0a:0a0a:ffff:ffff:ff04::
func (i *IP) String() (s string) {

	s = ""

	switch i.ver {

	case IPv4:
		for idx := range 4 {
			if s != "" {
				s += "."
			}
			s += strconv.Itoa(int(i.addr[idx]))
		}
		return

	case IPv6:
		for idx := range 8 {
			if s != "" {
				s += ":"
			}
			s += fmt.Sprintf("%02X%02X", i.addr[idx*2], i.addr[idx*2+1])
		}
		s = strings.ToLower(
			strings.ReplaceAll(s, ":0000", "::"),
		)

		return

	default:
		panic(ErrUncompVersion)
	}
}

func (i *IP) StringFull() (s string) {
	return fmt.Sprintf("%s/%v", i.String(), i.mask)
}

func (i *IP) safeSubnetSet(value byte) (ok bool) {
	if 0 > value {
		return
	}
	if (value > 32 && i.ver == IPv4) || (value > 128 && i.ver == IPv6) {
		return
	}
	i.mask = value
	return true
}

// hexStringToDecimal converts hex string to byte slice with fixed length
func hexStringToDecimal(hexStr string, len int) ([]byte, error) {
	decimals := make([]byte, len)

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return decimals, err
	}

	for i, b := range bytes {
		decimals[i] = b
	}

	return decimals, nil
}
