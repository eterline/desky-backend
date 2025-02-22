package configuration

import "github.com/google/uuid"

var FileName string = "settings.yaml"

var DefaultParameters = &Configuration{
	DevelopEnv: false,

	Server: Server{

		Address: ServerAddr{
			IP:   "0.0.0.0",
			Port: 3000,
		},

		SSL: ServerSSL{
			TLS:      false,
			CertFile: "",
			KeyFile:  "",
		},
	},

	Agent: AgentOptions{
		UUID:       uuid.NewString(),
		DefaultQoS: 1,
		Username:   "user",
		Password:   "user",
		Server: AgentServer{
			ConnectTimeout: "30s",
			Protocol:       "tcp",
			Host:           "localhost",
			Port:           1883,
		},
	},
}
