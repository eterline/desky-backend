package configuration

import "github.com/eterline/desky-backend/pkg/broker"

// ============================= Main app config struct =============================
type Configuration struct {
	DevelopEnv bool   `yaml:"dev-env" validate:"boolean"`
	DB         DB     `yaml:"DB"`
	Server     Server `yaml:"HTTP-Server" validate:"required"`
	Agent      AgentOptions
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

// ============================= DeskyAgent config struct =============================

type AgentOptions struct {
	UUID     string `yaml:"mqtt-id" validate:"required,uuid"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Server   AgentServer
}

type AgentServer struct {
	ConnectTimeout string             `yaml:"connect-timeout" validate:"required"`
	DefaultQoS     broker.QoSValue    `yaml:"default-qos" validate:"required,oneof=0 1 2 3"`
	Protocol       broker.SenderProto `yaml:"proto" validate:"required,oneof=ws ssl tcp"`
	Host           string             `yaml:"host" validate:"required,hostname"`
	Port           uint16             `yaml:"port" validate:"required,port"`
}

// Services config struct =============================

// -----Can be used with api pooling-----

// func (a DeskyAgent) ValueAPI() string {
// 	return a.API
// }
// func (a DeskyAgent) ValueToken() string {
// 	return a.Token
// }

// func (a DeskyAgentList) ToRequestList() []agentmon.AgentRequest {
// 	agentRequests := make([]agentmon.AgentRequest, len(a))

// 	for i, agent := range a {
// 		agentRequests[i] = agent
// 	}

// 	return agentRequests
// }

// ============================= Db config struct =============================

type (
	DB struct {
		File string `yaml:"file"`
		Sync bool   `yaml:"sync"`
	}
)
