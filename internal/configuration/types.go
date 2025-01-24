package configuration

// ============================= Main app config struct =============================
type Configuration struct {
	DevelopEnv bool               `yaml:"dev-env" validate:"boolean"`
	Server     HTTPServer         `yaml:"HTTP-Server" validate:"required"`
	Logs       Logging            `yaml:"Logs" validate:"required"`
	Services   ServicesParameters `yaml:"Services"`
}

// Server config struct =============================
type (
	HTTPServer struct {
		SSL     SSLParameters `yaml:"SSL"`
		Address Addr          `yaml:"Address" validate:"required"`
	}

	Addr struct {
		IP   string `yaml:"listen" validate:"required,ip"`
		Port uint16 `yaml:"port" validate:"required,port"`
	}

	SSLParameters struct {
		TLS      bool   `yaml:"tls-mode" validate:"boolean"`
		CertFile string `yaml:"cert-file" validate:"required"`
		KeyFile  string `yaml:"key-file" validate:"required"`
	}
)

// Logging config struct =============================
type (
	Logging struct {
		Enabled bool   `yaml:"enabled" validate:"boolean"`
		Level   int    `yaml:"level"`
		Path    string `yaml:"path"`
	}
)

// ============================= Services config struct =============================

type ServicesParameters struct {
	PVE []PVEInstance `yaml:"ProxmoxVE"`
}

// PVE config struct =============================
type (
	PVEInstance struct {
		Node     string `yaml:"node"`
		API      string `yaml:"api-url"`
		Username string `yaml:"username"`
		Secret   string `yaml:"secret"`
	}
)
