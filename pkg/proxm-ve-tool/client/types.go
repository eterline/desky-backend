package client

type (
	ProxmoxAuthResponse struct {
		Data Data `json:"data"`
	}

	Dc struct {
		SDNAllocate int `json:"SDN.Allocate"`
		SysModify   int `json:"Sys.Modify"`
		SDNAudit    int `json:"SDN.Audit"`
		SysAudit    int `json:"Sys.Audit"`
		SDNUse      int `json:"SDN.Use"`
	}
	Mapping struct {
		MappingAudit      int `json:"Mapping.Audit"`
		PermissionsModify int `json:"Permissions.Modify"`
		MappingModify     int `json:"Mapping.Modify"`
		MappingUse        int `json:"Mapping.Use"`
	}
	Vms struct {
		VMConfigDisk       int `json:"VM.Config.Disk"`
		VMConfigMemory     int `json:"VM.Config.Memory"`
		VMConfigCloudinit  int `json:"VM.Config.Cloudinit"`
		VMPowerMgmt        int `json:"VM.PowerMgmt"`
		VMMonitor          int `json:"VM.Monitor"`
		PermissionsModify  int `json:"Permissions.Modify"`
		VMMigrate          int `json:"VM.Migrate"`
		VMSnapshot         int `json:"VM.Snapshot"`
		VMBackup           int `json:"VM.Backup"`
		VMConfigNetwork    int `json:"VM.Config.Network"`
		VMClone            int `json:"VM.Clone"`
		VMConsole          int `json:"VM.Console"`
		VMConfigCPU        int `json:"VM.Config.CPU"`
		VMConfigOptions    int `json:"VM.Config.Options"`
		VMConfigCDROM      int `json:"VM.Config.CDROM"`
		VMSnapshotRollback int `json:"VM.Snapshot.Rollback"`
		VMConfigHWType     int `json:"VM.Config.HWType"`
		VMAllocate         int `json:"VM.Allocate"`
		VMAudit            int `json:"VM.Audit"`
	}

	Access struct {
		PermissionsModify int `json:"Permissions.Modify"`
		UserModify        int `json:"User.Modify"`
		GroupAllocate     int `json:"Group.Allocate"`
	}

	Nodes struct {
		SysModify         int `json:"Sys.Modify"`
		SysAudit          int `json:"Sys.Audit"`
		SysConsole        int `json:"Sys.Console"`
		SysIncoming       int `json:"Sys.Incoming"`
		PermissionsModify int `json:"Permissions.Modify"`
		SysSyslog         int `json:"Sys.Syslog"`
		SysPowerMgmt      int `json:"Sys.PowerMgmt"`
		SysAccessNetwork  int `json:"Sys.AccessNetwork"`
	}

	Storage struct {
		DatastoreAudit            int `json:"Datastore.Audit"`
		DatastoreAllocateTemplate int `json:"Datastore.AllocateTemplate"`
		DatastoreAllocate         int `json:"Datastore.Allocate"`
		PermissionsModify         int `json:"Permissions.Modify"`
		DatastoreAllocateSpace    int `json:"Datastore.AllocateSpace"`
	}

	Sdn struct {
		SDNUse            int `json:"SDN.Use"`
		PermissionsModify int `json:"Permissions.Modify"`
		SDNAllocate       int `json:"SDN.Allocate"`
		SDNAudit          int `json:"SDN.Audit"`
	}

	Cap struct {
		Dc      Dc      `json:"dc"`
		Mapping Mapping `json:"mapping"`
		Vms     Vms     `json:"vms"`
		Access  Access  `json:"access"`
		Nodes   Nodes   `json:"nodes"`
		Storage Storage `json:"storage"`
		Sdn     Sdn     `json:"sdn"`
	}

	Data struct {
		Ticket              string `json:"ticket"`
		Username            string `json:"username"`
		Cap                 Cap    `json:"cap"`
		CSRFPreventionToken string `json:"CSRFPreventionToken"`
	}
)
