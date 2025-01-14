package configuration

import (
	"encoding/json"
	"os"
)

var defaultConfig = &Configuration{

	DevMode: false,

	Auth: AuthParameters{
		Enabled:  true,
		Bcrypt:   false,
		Username: "admin",
		Password: "admin",
	},

	Server: ServerConfig{

		JWTSecret: "QO8hzZPlWeNYTg0UeO7wVaddySdUwo6HYSPQbemjxRiCCVj6",

		TLS: ServerTLSConfig{
			Enabled:     false,
			Key:         "",
			Certificate: "",
		},
		Connection: ServerConnectionConfig{
			Addr:     "0.0.0.0",
			Hostname: "",
			Port:     3000,
		},
	},

	Services: map[string]ServiceParams{
		"default-service-0": ServiceParams{
			IsActive: true,
			ApiURL:   "http://example-api.lan/api",
			UseKeys:  false,
			Username: "admin",
			Password: "admin",
		},
		"default-service-1": ServiceParams{
			IsActive: true,
			ApiURL:   "http://example-api.lan/api",
			UseKeys:  true,
			Key:      "NNNUEiomn39488f9945f894hfv8u4nuv034v89v3u89v508huig",
		},
	},

	Proxmox: []ProxmoxCredentials{
		{
			SSLCheck: false,
			ApiURL:   "http://proxmox.lan:8006",
			Username: "root",
			Password: "root",
		},
	},
}

func GenerateFile(out string) error {
	file, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()

	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(defaultConfig, "", "   ")
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}
