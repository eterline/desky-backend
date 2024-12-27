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
		Cores     int64   `json:"cores"`
		Frequency int64   `json:"frequency"`
	}
)
