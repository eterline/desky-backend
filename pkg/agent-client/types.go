package agentclient

import (
	"fmt"
	"net/url"
)

var ExportersMap map[string]any = map[string]any{
	"host":        Host{},
	"cpu":         CPU{},
	"ram":         RAM{},
	"load":        Load{},
	"temperature": SensorList{},
	"ports":       Ports{},
	"partitions":  PartitionList{},
}

type HostInfo struct {
	HostID   string `json:"HostID"`
	Hostname string `json:"Hostname"`
}

type DeskyAgent struct {
	Url  *url.URL
	Info *HostInfo

	token keyBerearer
}

type keyBerearer string

func (key keyBerearer) Berearer() string {
	return fmt.Sprintf("Berearer %s", key)
}

// Represent objectsAPI

type (
	Host struct {
		Name         string  `json:"hostname"`
		Uptime       float64 `json:"uptime"`
		OS           string  `json:"os"`
		ProcessCount uint64  `json:"processes"`
		VirtSystem   string  `json:"hypervisor"`
	}
)

type (
	CPU struct {
		Name        string    `json:"name"`
		Model       string    `json:"model"`
		CoreCount   uint64    `json:"core-count"`
		ThreadCount uint64    `json:"thread-count"`
		Cores       []CpuCore `json:"cores"`
		Cache       int32     `json:"cache"`
		Load        float64   `json:"load"`
	}

	CpuCore struct {
		ID      string  `json:"id"`
		FreqMhz float64 `json:"frequency"`
	}

	Load struct {
		Load1  float64 `json:"load-1"`
		Load5  float64 `json:"load-5"`
		Load15 float64 `json:"load-15"`
	}
)

type (
	RAM struct {
		Total      uint64  `json:"total"`
		Used       uint64  `json:"used"`
		Avail      uint64  `json:"available"`
		UsePercent float64 `json:"use"`
	}
)
type (
	SensorList []Sensor

	Sensor struct {
		Key     string  `json:"key"`
		Current float64 `json:"current"`
		Max     float64 `json:"max"`
	}
)

type (
	Ports []Interface

	Interface struct {
		MTU          int    `json:"mtu"`  // maximum transmission unit
		Name         string `json:"name"` // e.g., "en0", "lo0", "eth0.100"
		HardwareAddr string `json:"mac"`  // IEEE MAC-48, EUI-48 and EUI-64 form
	}

	InterfaceAddr struct {
		Addr string `json:"addr"`
	}

	Connect struct {
		Type   uint32 `json:"type"`
		Status string `json:"status"`
		Laddr  string `json:"local-addr"`
		Raddr  string `json:"remote-addr"`
		Pid    int32  `json:"pid"`
	}
)

type (
	PartitionList []Partition

	Partition struct {
		Device      string  `json:"device"`
		FS          string  `json:"fs"`
		Total       uint64  `json:"total"`
		Free        uint64  `json:"free"`
		Used        uint64  `json:"used"`
		UsedPercent float64 `json:"used-percent"`
	}
)
