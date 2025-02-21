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
		UUID:     uuid.NewString(),
		Username: "user",
		Password: "user",
		Server: AgentServer{
			ConnectTimeout: "30s",
			DefaultQoS:     1,
			Protocol:       "tcp",
			Host:           "localhost",
			Port:           1883,
		},

		// -----------uses with api pooling form-----------
		// DeskyAgent: []DeskyAgent{
		// 	{
		// 		API:   "http://host.lan/api",
		// 		Token: "insert-example-24digit-token",
		// 	},
		// },
	},
}
