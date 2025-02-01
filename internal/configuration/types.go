package configuration

// ============================= Main app config struct =============================
type Configuration struct {
	DevelopEnv bool     `yaml:"dev-env" validate:"boolean"`
	DB         DB       `yaml:"DB"`
	Server     Server   `yaml:"HTTP-Server" validate:"required"`
	Services   Services `yaml:"Services"`
}

// Server config struct =============================
type (
	Server struct {
		Name    string     `yaml:Name"`
		Address ServerAddr `yaml:"Address" validate:"required"`
		SSL     ServerSSL  `yaml:"SSL"`
	}

	ServerAddr struct {
		IP   string `yaml:"listen" validate:"required,ip"`
		Port uint16 `yaml:"port" validate:"required,port"`
	}

	ServerSSL struct {
		TLS      bool   `yaml:"tls-mode" validate:"boolean"`
		CertFile string `yaml:"cert-file" validate:"required"`
		KeyFile  string `yaml:"key-file" validate:"required"`
	}
)

// Logging config struct =============================

// ============================= Services config struct =============================

type Services struct {
	Proxmox    []ProxmoxInstance `yaml:"Proxmox"`
	Docker     []DockerInstance  `yaml:"Docker"`
	DeskyAgent []DeskyAgent      `yaml:"DeskyAgent"`
}

// Services config struct =============================
type (
	ProxmoxInstance struct {
		Node     string `yaml:"node"`
		API      string `yaml:"api-url"`
		Username string `yaml:"username"`
		Secret   string `yaml:"secret"`
	}

	DockerInstance struct {
		Name string `yaml:"name"`
		API  string `yaml:"api"`
	}

	DeskyAgent struct {
		API   string `yaml:"api"`
		Token string `yaml:"token"`
	}
)

// ============================= Db config struct =============================

type (
	DB struct {
		File string `yaml:"file"`
		Sync bool   `yaml:"sync"`
	}
)
