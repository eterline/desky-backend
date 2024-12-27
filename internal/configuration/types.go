package configuration

// ============================= Main app config field =============================
type Configuration struct {
	DevMode  bool           `json:"Development" yaml:"Development" validate:"boolean"`
	Server   ServerConfig   `json:"Server" yaml:"Server" validate:"required"`
	Services ServiceMap     `json:"Services" yaml:"Services" validate:"required"`
	Proxmox  ProxmoxService `json:"Proxmox" yaml:"Proxmox" validate:"required"`
	// Storage

	_ struct{}
}

// Server config field =============================

type (
	// Server configrurations
	ServerConfig struct {
		TLS        ServerTLSConfig        `json:"SSL-Mode" yaml:"SSL-Mode"`
		Connection ServerConnectionConfig `json:"Connection" yaml:"Connection" validate:"required"`

		_ struct{}
	}

	// HTTPS Presets for web server
	ServerTLSConfig struct {
		Enabled     bool   `json:"active" yaml:"active" validate:"boolean"`
		Key         string `json:"key-file" yaml:"key-file" validate:"omitempty,filepath"`
		Certificate string `json:"cert-file" yaml:"cert-file" validate:"omitempty,filepath"`

		_ struct{}
	}

	// Web server parameters.
	// Listening address of HostName:Port
	ServerConnectionConfig struct {
		Addr     string `json:"address" yaml:"address" validate:"required,ip"`
		Hostname string `json:"host-name" yaml:"host-name"`
		Port     uint16 `json:"port" yaml:"port" validate:"required,port"`

		_ struct{}
	}
)

// Services config field =============================

type (
	// Abstract service paramaters.
	// That can be used with APIkey-like using or typical credentials
	ServiceParams struct {
		IsActive bool   `json:"active" yaml:"active" validate:"boolean"`
		ApiURL   string `json:"api-url" yaml:"api-url" validate:"required, url"`
		UseKeys  bool   `json:"use-api-key" yaml:"use-api-key" validate:"boolean"`
		Key      string `json:"secret" yaml:"secret" validate:"required_if=UseKeys true"`
		Username string `json:"login" yaml:"login" validate:"required_if=UseKeys false"`
		Password string `json:"password" yaml:"password" validate:"required_if=UseKeys false"`

		_ struct{}
	}

	ServiceMap map[string]ServiceParams
)

type (
	ProxmoxService []ProxmoxCredentials

	ProxmoxCredentials struct {
		SSLCheck bool   `json:"ssl-check" yaml:"ssl-check" validate:"boolean"`
		ApiURL   string `json:"api-url" yaml:"api-url" validate:"required, url"`
		Username string `json:"login" yaml:"login"`
		Password string `json:"password" yaml:"password"`

		_ struct{}
	}
)

// Storage config field =============================
