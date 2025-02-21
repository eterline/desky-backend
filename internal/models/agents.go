package models

var ExporterList = []string{"host", "cpu", "ram", "load", "temperature", "ports", "partitions"}

type SessionCredentials struct {
	Hostname string `json:"hostname"`
	ID       string `json:"id"`
	URL      string `json:"url"`
}

type FetchedResponse struct {
	ID   string         `json:"id"`
	Data map[string]any `json:"data"`
}

type FetchedResponseSingle struct {
	ID   string `json:"id"`
	Data any    `json:"data"`
}

type CPUCore struct {
	ID        string  `json:"id"`
	Frequency float64 `json:"frequency"`
}

type CPU struct {
	Name        string    `json:"name"`
	Model       string    `json:"model"`
	Cache       int       `json:"cache"`
	CoreCount   int       `json:"core-count"`
	ThreadCount int       `json:"thread-count"`
	Load        float64   `json:"load"`
	Cores       []CPUCore `json:"cores"`
}

type Host struct {
	Hostname   string `json:"hostname"`
	Hypervisor string `json:"hypervisor"`
	OS         string `json:"os"`
	Processes  int    `json:"processes"`
	Uptime     int    `json:"uptime"`
}

type Load struct {
	Load1  float64 `json:"load-1"`
	Load5  float64 `json:"load-5"`
	Load15 float64 `json:"load-15"`
}

type Partition struct {
	Device      string  `json:"device"`
	FS          string  `json:"fs"`
	Total       int     `json:"total"`
	Free        int     `json:"free"`
	Used        int     `json:"used"`
	UsedPercent float64 `json:"used-percent"`
}

type Port struct {
	MAC  string `json:"mac"`
	MTU  int    `json:"mtu"`
	Name string `json:"name"`
}

type RAM struct {
	Total     int     `json:"total"`
	Available int     `json:"available"`
	Used      int     `json:"used"`
	Use       float64 `json:"use"`
}

type Temperature struct {
	Key     string  `json:"key"`
	Current float64 `json:"current"`
	Max     float64 `json:"max"`
}

type AgentStatsObject struct {
	CPU         *CPU           `json:"cpu"`
	Host        *Host          `json:"host"`
	Load        *Load          `json:"load"`
	Partitions  *[]Partition   `json:"partitions"`
	Ports       *[]Port        `json:"ports"`
	RAM         *RAM           `json:"ram"`
	Temperature *[]Temperature `json:"temperature"`
}
