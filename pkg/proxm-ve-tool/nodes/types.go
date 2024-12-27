package nodes

type (
	NodeList struct {
		Data []NodeUnit `json:"data"`
		Code int
	}

	NodeUnit struct {
		CPU            float64 `json:"cpu"`
		Status         string  `json:"status"`
		Maxcpu         int     `json:"maxcpu"`
		Mem            int64   `json:"mem"`
		Node           string  `json:"node"`
		Disk           int64   `json:"disk"`
		ID             string  `json:"id"`
		SslFingerprint string  `json:"ssl_fingerprint"`
		Uptime         int     `json:"uptime"`
		Level          string  `json:"level"`
		Maxmem         int64   `json:"maxmem"`
		Maxdisk        int64   `json:"maxdisk"`
		Type           string  `json:"type"`
	}
)

// node options structs =============

// /status
type (
	NodeStatus struct {
		Data NodeStatusData `json:"data"`
		Code int
	}
	Memory struct {
		Used  int64 `json:"used"`
		Free  int64 `json:"free"`
		Total int64 `json:"total"`
	}
	Cpuinfo struct {
		Hvm     string `json:"hvm"`
		Mhz     string `json:"mhz"`
		Model   string `json:"model"`
		Sockets int    `json:"sockets"`
		UserHz  int    `json:"user_hz"`
		Cpus    int    `json:"cpus"`
		Cores   int    `json:"cores"`
		Flags   string `json:"flags"`
	}
	Rootfs struct {
		Used  int64 `json:"used"`
		Total int64 `json:"total"`
		Avail int64 `json:"avail"`
		Free  int64 `json:"free"`
	}
	CurrentKernel struct {
		Machine string `json:"machine"`
		Sysname string `json:"sysname"`
		Release string `json:"release"`
		Version string `json:"version"`
	}
	Swap struct {
		Used  int   `json:"used"`
		Total int64 `json:"total"`
		Free  int64 `json:"free"`
	}
	BootInfo struct {
		Secureboot int    `json:"secureboot"`
		Mode       string `json:"mode"`
	}
	Ksm struct {
		Shared int `json:"shared"`
	}
	NodeStatusData struct {
		Memory        Memory        `json:"memory"`
		Loadavg       []string      `json:"loadavg"`
		Wait          float64       `json:"wait"`
		Uptime        int           `json:"uptime"`
		Kversion      string        `json:"kversion"`
		Cpuinfo       Cpuinfo       `json:"cpuinfo"`
		Rootfs        Rootfs        `json:"rootfs"`
		Idle          int           `json:"idle"`
		CurrentKernel CurrentKernel `json:"current-kernel"`
		Swap          Swap          `json:"swap"`
		BootInfo      BootInfo      `json:"boot-info"`
		Ksm           Ksm           `json:"ksm"`
		Pveversion    string        `json:"pveversion"`
		CPU           float64       `json:"cpu"`
	}
)

// /hosts
type (
	HostsFile struct {
		Data HostsFileData `json:"data"`
		Code int
	}

	HostsFileData struct {
		Digest string `json:"digest"`
		Data   string `json:"data"`
	}
)

// /dns
type (
	DNS struct {
		Data DNSData `json:"data"`
		Code int
	}

	DNSData struct {
		Search string `json:"search"`
		DNS1   string `json:"dns1"`
	}
)

// /aplinfo
type (
	AplInfo struct {
		Data []AplInfoPage `json:"data"`
		Code int
	}

	AplInfoPage struct {
		Section      string `json:"section"`
		Sha512Sum    string `json:"sha512sum"`
		Location     string `json:"location"`
		Template     string `json:"template"`
		Architecture string `json:"architecture"`
		Md5Sum       string `json:"md5sum"`
		Maintainer   string `json:"maintainer"`
		Source       string `json:"source"`
		Headline     string `json:"headline"`
		Version      string `json:"version"`
		Package      string `json:"package"`
		Description  string `json:"description"`
		Infopage     string `json:"infopage"`
		Type         string `json:"type"`
		Os           string `json:"os"`
	}
)

// /lxc
type (
	LXCList struct {
		Data []LXCDataUnit `json:"data"`
		Code int
	}

	LXCDataUnit struct {
		Tags      string  `json:"tags"`
		Name      string  `json:"name"`
		Cpus      int     `json:"cpus"`
		Swap      int     `json:"swap"`
		Uptime    int     `json:"uptime"`
		Maxswap   int     `json:"maxswap"`
		Mem       int     `json:"mem"`
		Maxdisk   int64   `json:"maxdisk"`
		Netin     int     `json:"netin"`
		Diskread  int     `json:"diskread"`
		Pid       int     `json:"pid"`
		Netout    int     `json:"netout"`
		Disk      int     `json:"disk"`
		Status    string  `json:"status"`
		CPU       float64 `json:"cpu"`
		Vmid      int     `json:"vmid"`
		Diskwrite int     `json:"diskwrite"`
		Type      string  `json:"type"`
		Maxmem    int     `json:"maxmem"`
	}
)

// /qemu
type (
	VMList struct {
		Data []VMDataUnit `json:"data"`
		Code int
	}

	VMDataUnit struct {
		Tags      string  `json:"tags"`
		Name      string  `json:"name"`
		Cpus      int     `json:"cpus"`
		Uptime    int     `json:"uptime"`
		Mem       int     `json:"mem"`
		Maxdisk   int64   `json:"maxdisk"`
		Netin     int     `json:"netin"`
		Diskread  int     `json:"diskread"`
		Netout    int     `json:"netout"`
		Pid       int     `json:"pid"`
		Disk      int     `json:"disk"`
		Status    string  `json:"status"`
		CPU       float64 `json:"cpu"`
		Diskwrite int     `json:"diskwrite"`
		Vmid      int     `json:"vmid"`
		Maxmem    int64   `json:"maxmem"`
	}
)
