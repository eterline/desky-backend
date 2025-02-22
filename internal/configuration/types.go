package configuration

import "github.com/eterline/desky-backend/pkg/broker"

// ============================= Main app config struct =============================
type Configuration struct {
	DevelopEnv bool         `yaml:"dev-env" validate:"boolean"`
	DB         DB           `yaml:"DB"`
	Server     Server       `yaml:"HTTP-Server" validate:"required"`
	Agent      AgentOptions `yaml:AgentMQTT"`
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
		CertFile string `yaml:"cert-file" validate:"required_if=TLS true"`
		KeyFile  string `yaml:"key-file" validate:"required_if=TLS  true"`
	}
)

// Logging config struct =============================

// ============================= DeskyAgent config struct =============================

type AgentOptions struct {
	UUID       string          `yaml:"mqtt-uuid" validate:"required,uuid"`
	DefaultQoS broker.QoSValue `yaml:"default-qos" validate:"oneof=0 1 2 3"`
	Username   string          `yaml:"Username"`
	Password   string          `yaml:"Password"`
	Server     AgentServer     `yaml:"Server"`
}

type AgentServer struct {
	Protocol       broker.SenderProto `yaml:"proto" validate:"required,oneof=ws ssl tcp"`
	Host           string             `yaml:"host" validate:"required,hostname"`
	Port           uint16             `yaml:"port" validate:"required,port"`
	ConnectTimeout string             `yaml:"connect-timeout" validate:"required"`
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
