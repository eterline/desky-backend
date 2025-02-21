package configuration

import (
	"fmt"
	"time"

	"github.com/eterline/desky-backend/pkg/iptool"
)

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

func (c *Configuration) SSL() ServerSSL {
	return c.Server.SSL
}

func (c *Configuration) URLString() string {

	var proto string

	if c.SSL().TLS {
		proto = "https"
	} else {
		proto = "http"
	}

	return fmt.Sprintf("%s://%s", proto, c.ServerSocket())
}

func (c *Configuration) MQTTConnTimeout() time.Duration {

	tm, err := time.ParseDuration(c.Agent.Server.ConnectTimeout)
	if err != nil {
		return 30 * time.Second
	}

	return tm
}
