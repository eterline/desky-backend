package ve

type (
	PVENodeStatus struct {
		Name    string      `json:"name"`
		AVGLoad AVGLoadData `json:"load"`
		FS      FSData      `json:"fs"`
		RAM     RAMData     `json:"ram"`
		CPU     CPUData     `json:"cpu"`
		Uptime  int64       `json:"uptime"`
		Kernel  string      `json:"kernel"`
	}

	AVGLoadData [3]string

	FSData struct {
		Used  int64 `json:"used"`
		Total int64 `json:"total"`
	}

	RAMData struct {
		Used  int64 `json:"used"`
		Total int64 `json:"total"`
	}

	CPUData struct {
		Load      float64 `json:"load"`
		Model     string  `json:"model"`
		Cores     int     `json:"cores"`
		Frequency string  `json:"frequency"`
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
