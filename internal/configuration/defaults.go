package configuration

var DefaultParameters = &Configuration{
	DevelopEnv: false,

	Server: HTTPServer{

		Address: Addr{
			IP:   "0.0.0.0",
			Port: 3000,
		},

		SSL: SSLParameters{
			TLS:      false,
			CertFile: "",
			KeyFile:  "",
		},
	},

	Logs: Logging{
		Enabled: true,
		Level:   0,
		Path:    "./logs",
	},

	Services: ServicesParameters{
		PVE: []PVEInstance{
			PVEInstance{
				Node:     "node-1",
				API:      "https://node-1:8006/api/json",
				Username: "root",
				Secret:   "example_key",
			},
			PVEInstance{
				Node:     "node-2",
				API:      "https://node-2:8006/api/json",
				Username: "root",
				Secret:   "example_key",
			},
		},
	},
}
