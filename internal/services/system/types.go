package system

type (
	SystemdUnit struct {
		UnitFile string `json:"unit_file"`
		Status   string `json:"state"`
		Preset   string `json:"preset"`
	}
)

type (
	HostInfo struct {
		Name         string         `json:"hostname"`
		Uptime       UptimeDuration `json:"uptime"`
		OS           string         `json:"os"`
		ProcessCount uint64         `json:"procs"`
		VirtSystem   string         `json:"virt"`
		Addrs        AddrsList      `json:"addrs"`
	}

	AddrsList      []string
	UptimeDuration float64
)

type (
	RAMInfo struct {
		Total      uint64 `json:"total"`
		Used       uint64 `json:"used"`
		Avail      uint64 `json:"available"`
		UsePercent uint8  `json:"use"`
	}
)

type (
	CPUInfo struct {
		Name        string    `json:"name"`
		Model       string    `json:"model"`
		CoreCount   uint64    `json:"coreCount"`
		ThreadCount uint64    `json:"threadCount"`
		Cores       []CpuCore `json:"cores"`
		Cache       int32     `json:"cache"`
	}

	CpuCore struct {
		ID      string  `json:"id"`
		FreqMhz float64 `json:"frequency"`
	}
)

type (
	SensorInfo struct {
		Key     string  `json:"key"`
		Current float64 `json:"current"`
		Max     float64 `json:"max"`
	}
)

type (
	RequestCLI struct {
		Command string `json:"command"`
	}
	ResponseCLI struct {
		Command string `json:"command"`
		Output  string `json:"output"`
	}
)
