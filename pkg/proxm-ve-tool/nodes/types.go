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

type (
	AptVersions struct {
		Data []AptVersionUnit `json:"data"`
		Code int
	}

	AptVersionUnit struct {
		Priority       string `json:"Priority"`
		Version        string `json:"Version"`
		Arch           string `json:"Arch"`
		Description    string `json:"Description"`
		CurrentState   string `json:"CurrentState"`
		OldVersion     string `json:"OldVersion"`
		ManagerVersion string `json:"ManagerVersion"`
		Origin         string `json:"Origin"`
		Title          string `json:"Title"`
		Section        string `json:"Section"`
		Package        string `json:"Package"`
	}
)

type (
	AptUpdates struct {
		Data []AptUpdateUnit `json:"data"`
		Code int
	}

	AptUpdateUnit struct {
		Description string `json:"Description,omitempty"`
		Origin      string `json:"Origin,omitempty"`
		OldVersion  string `json:"OldVersion,omitempty"`
		Version     string `json:"Version,omitempty"`
		Section     string `json:"Section,omitempty"`
		Priority    string `json:"Priority,omitempty"`
		Arch        string `json:"Arch,omitempty"`
		Package     string `json:"Package,omitempty"`
		Title       string `json:"Title,omitempty"`
	}
)

type (
	AptUpgrade struct {
		Data []string `json:"data"`
		Code int
	}
)

type (
	DisksInfo struct {
		Data []DiskUnit `json:"data"`
		Code int
	}

	DiskUnit struct {
		Wwn       string `json:"wwn"`
		Used      string `json:"used"`
		Devpath   string `json:"devpath"`
		Osdid     int    `json:"osdid"`
		Model     string `json:"model"`
		OsdidList any    `json:"osdid-list"`
		Serial    string `json:"serial"`
		Type      string `json:"type"`
		Vendor    string `json:"vendor"`
		ByIDLink  string `json:"by_id_link"`
		Size      int64  `json:"size"`
		Rpm       int    `json:"rpm"`
		Wearout   int    `json:"wearout"`
		Health    string `json:"health"`
		Gpt       int    `json:"gpt"`
		pn        *ProxmoxNode
	}

	DiskSmartData struct {
		Data DiskSmart `json:"data"`
		Code int
	}

	DiskSmart struct {
		Text    string `json:"text"`
		Wearout int    `json:"wearout"`
		Type    string `json:"type"`
		Health  string `json:"health"`
	}

	Smart struct {
		Info    SmartInfo `json:"info"`
		Wearout int       `json:"wearout"`
		Type    string    `json:"type"`
		Health  string    `json:"health"`
	}

	SmartInfo map[string]string
)
