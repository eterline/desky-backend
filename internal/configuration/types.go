package configuration

import agentmon "github.com/eterline/desky-backend/internal/services/agent-mon"

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
	Proxmox    ProxmoxInstanceList `yaml:"Proxmox"`
	Docker     DockerInstanceList  `yaml:"Docker"`
	DeskyAgent DeskyAgentList      `yaml:"DeskyAgent"`
}

// Services config struct =============================
type (
	ProxmoxInstanceList []ProxmoxInstance
	ProxmoxInstance     struct {
		Node     string `yaml:"node"`
		API      string `yaml:"api-url"`
		Username string `yaml:"username"`
		Secret   string `yaml:"secret"`
	}

	DockerInstanceList []DockerInstance
	DockerInstance     struct {
		Name string `yaml:"name"`
		API  string `yaml:"api"`
	}

	DeskyAgentList []DeskyAgent
	DeskyAgent     struct {
		API   string `yaml:"api"`
		Token string `yaml:"token"`
	}
)

func (a DeskyAgent) ValueAPI() string {
	return a.API
}
func (a DeskyAgent) ValueToken() string {
	return a.Token
}

func (a DeskyAgentList) ToRequestList() []agentmon.AgentRequest {
	agentRequests := make([]agentmon.AgentRequest, len(a))

	for i, agent := range a {
		agentRequests[i] = agent
	}

	return agentRequests
}

// ============================= Db config struct =============================

type (
	DB struct {
		File string `yaml:"file"`
		Sync bool   `yaml:"sync"`
	}
)
