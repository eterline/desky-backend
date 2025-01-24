package configuration

import "github.com/eterline/desky-backend/pkg/iptool"

func (c *Configuration) ServerSocket() string {

	ip, ok := iptool.ParseIPv4(c.Server.Address.IP)
	if !ok {
		ip = iptool.ZeroIP
	}

	sock, ok := iptool.NewSocket(ip, int(c.Server.Address.Port))
	if !ok {
		sock, _ = iptool.NewSocket(ip, 3000)
	}

	return sock.String()
}

func (c *Configuration) SSL() SSLParameters {
	return c.Server.SSL
}
