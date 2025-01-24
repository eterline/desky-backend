package models

type (
	PVENodeStatus struct {
		Name    string      `json:"name" example"micro-ve"`
		AVGLoad AVGLoadData `json:"load"`
		FS      FSData      `json:"fs"`
		RAM     RAMData     `json:"ram"`
		CPU     CPUData     `json:"cpu"`
		Uptime  int64       `json:"uptime" example"512"`
		Kernel  string      `json:"kernel" example"pve"`
	}

	AVGLoadData [3]string

	FSData struct {
		Used  int64 `json:"used" example"13255"`
		Total int64 `json:"total" example"13255"`
	}

	RAMData struct {
		Used  int64 `json:"used" example:"3220"`
		Total int64 `json:"total" example:"7680"`
	}

	CPUData struct {
		Load      float64 `json:"load" example:"13"`
		Model     string  `json:"model" example:"5"`
		Cores     int     `json:"cores" example:"6"`
		Frequency string  `json:"frequency" example:"4300Mhz"`
	}
)

type (
	DevicesList struct {
		LXC  []TypeDevice `json:"lxc"`
		QEMU []TypeDevice `json:"qemu"`
	}
)

type (
	TypeDevice struct {
		Status string `json:"status"`
		Name   string `json:"name"`
		Tags   string `json:"tags"`
		VMID   int    `json:"vmid"`
		CPUS   int    `json:"cpus"`
		NetIn  int64  `json:"netRX,omitempty"`
		NetOut int64  `json:"netTX,omitempty"`
		Uptime int64  `json:"uptime"`
		PID    int    `json:"pid"`
	}
)
