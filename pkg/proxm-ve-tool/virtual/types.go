package virtual

const (
	VirtTypeQEMU = "qemu"
	VirtTypeLXC  = "lxc"
)

type (
	QEMUStatus struct {
		Data QEMUStatusData `json:"data"`
		Code int
	}
	Ha struct {
		Managed int `json:"managed"`
	}
	Tap2000I0 struct {
		Netin  int `json:"netin"`
		Netout int `json:"netout"`
	}
	Nics struct {
		Tap2000I0 Tap2000I0 `json:"tap2000i0"`
	}
	ProxmoxSupport struct {
		PbsLibraryVersion       string `json:"pbs-library-version"`
		PbsDirtyBitmapMigration bool   `json:"pbs-dirty-bitmap-migration"`
		PbsMasterkey            bool   `json:"pbs-masterkey"`
		PbsDirtyBitmapSavevm    bool   `json:"pbs-dirty-bitmap-savevm"`
		BackupMaxWorkers        bool   `json:"backup-max-workers"`
		BackupFleecing          bool   `json:"backup-fleecing"`
		QueryBitmapInfo         bool   `json:"query-bitmap-info"`
		PbsDirtyBitmap          bool   `json:"pbs-dirty-bitmap"`
	}
	QEMUStatusData struct {
		RunningQemu    string         `json:"running-qemu"`
		Tags           string         `json:"tags"`
		Ha             Ha             `json:"ha"`
		Uptime         int            `json:"uptime"`
		Name           string         `json:"name"`
		Qmpstatus      string         `json:"qmpstatus"`
		Mem            int            `json:"mem"`
		Diskread       int            `json:"diskread"`
		Disk           int            `json:"disk"`
		Status         string         `json:"status"`
		CPU            float64        `json:"cpu"`
		Cpus           int            `json:"cpus"`
		Maxdisk        int64          `json:"maxdisk"`
		RunningMachine string         `json:"running-machine"`
		Netin          int            `json:"netin"`
		Nics           Nics           `json:"nics"`
		Pid            int            `json:"pid"`
		Netout         int            `json:"netout"`
		ProxmoxSupport ProxmoxSupport `json:"proxmox-support"`
		Diskwrite      int            `json:"diskwrite"`
		Vmid           int            `json:"vmid"`
		Maxmem         int64          `json:"maxmem"`
	}
)

type (
	LXCStatus struct {
		Data LXCStatusData `json:"data"`
		Code int
	}
	HaStatus struct {
		Managed int `json:"managed"`
	}
	LXCStatusData struct {
		Status    string   `json:"status"`
		Maxmem    int64    `json:"maxmem"`
		Diskwrite int      `json:"diskwrite"`
		Vmid      int      `json:"vmid"`
		Type      string   `json:"type"`
		CPU       int      `json:"cpu"`
		Disk      int      `json:"disk"`
		Netin     int      `json:"netin"`
		Maxdisk   int64    `json:"maxdisk"`
		Mem       int      `json:"mem"`
		Netout    int      `json:"netout"`
		Diskread  int      `json:"diskread"`
		Ha        HaStatus `json:"ha"`
		Tags      string   `json:"tags"`
		Maxswap   int      `json:"maxswap"`
		Uptime    int      `json:"uptime"`
		Name      string   `json:"name"`
		Cpus      int      `json:"cpus"`
		Swap      int      `json:"swap"`
	}
)
