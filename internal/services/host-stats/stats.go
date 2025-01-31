package hoststats

import (
	"fmt"
)

// TODO: Доделать работу с хостами

type HostInstance struct {
	Conn string
	key  []byte
}

type HostStatsService struct {
	HostStack []*HostInstance
}

func New() *HostStatsService {
	return new(HostStatsService)
}

func (hs *HostStatsService) RegisterHost(user, host string, port uint16, sshKey []byte) {

	if user == "" {
		user = "root"
	}

	i := &HostInstance{
		Conn: fmt.Sprintf("%s@%s:%v", user, host, port),
		key:  sshKey,
	}

	hs.HostStack = append(hs.HostStack, i)
}
