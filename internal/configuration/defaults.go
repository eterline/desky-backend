package configuration

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

	Services: Services{
		Proxmox: []ProxmoxInstance{
			{
				Node:     "node-1",
				API:      "https://node-1:8006/api/json",
				Username: "root",
				Secret:   "example_key",
			},
			{
				Node:     "node-2",
				API:      "https://node-2:8006/api/json",
				Username: "root",
				Secret:   "example_key",
			},
		},
		Docker: []DockerInstance{
			{
				Name: "docker-1",
				API:  "tcp://docker.lan",
			},
		},
		DeskyAgent: []DeskyAgent{
			{
				API:   "http://host.lan/api",
				Token: "insert-example-24digit-token",
			},
		},
	},
}
