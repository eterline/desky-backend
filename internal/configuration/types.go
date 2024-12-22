package configuration

// ============================= Main app config field =============================
type Configuration struct {
	DevMode  bool
	Server   ServerConfig
	Services ServiceMap
	// Storage

	_ struct{}
}

// Server config field =============================

type (
	// Server configrurations
	ServerConfig struct {
		TLS        ServerTLSConfig
		Connection ServerConnectionConfig

		_ struct{}
	}

	// HTTPS Presets for web server
	ServerTLSConfig struct {
		Enabled          bool
		Key, Certificate string

		_ struct{}
	}

	// Web server parameters.
	// Listening address of HostName:Port
	ServerConnectionConfig struct {
		Addr     string
		Hostname string
		Port     uint16

		_ struct{}
	}
)

// Services config field =============================

type (
	// Abstract service paramaters.
	// That can be used with APIkey-like using or typical credentials
	ServiceParams struct {
		IsActive           bool
		ApiURL             string
		UseKeys            bool
		Key                string
		Username, Password string

		_ struct{}
	}

	ServiceMap map[string]ServiceParams
)

// Storage config field =============================
