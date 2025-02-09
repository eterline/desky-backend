package models

// ====================================================

type SSHtypeOS string

const (
	Nothing SSHtypeOS = "none"
	Windows SSHtypeOS = "windows"
	Linux   SSHtypeOS = "linux"
	BSD     SSHtypeOS = "bsd"
)

func StringSSHtypeOS(value string) SSHtypeOS {

	switch SSHtypeOS(value) {
	case Windows:
		return Windows
	case Linux:
		return Linux
	case BSD:
		return BSD
	default:
		return Nothing
	}
}

// ====================================================

type RequestFormSSH struct {
	System string `json:"os" validate:"required"`
	Port   uint16 `json:"port" validate:"port"`
	Host   string `json:"host" validate:"required"`

	User          string `json:"user" validate:"required"`
	PrivateKeyUse bool   `json:"private-key-use" validate:"boolean"`

	Password   string `json:"password" validate:"required_if=PrivateKeyUse false"`
	PrivateKey string `json:"private-key" validate:"required_if=PrivateKeyUse true"`
}

type ResponseCreateSSH struct {
	PrivateKeyUse bool   `json:"private-key-use"`
	Target        string `json:"target"`
}

type SSHInstanceObject struct {
	ID            int    `json:"id"`
	HostString    string `json:"host"`
	PrivateKeyUse bool   `json:"private-key-use"`
}

type SSHTestObject struct {
	ID        int  `json:"id"`
	Available bool `json:"available"`
}

// ====================================================

type SSHSessionResponse struct {
	Host string `json:"host"`
	User string `json:"user"`

	Command  string `json:"command"`
	Response string `json:"response"`
	Closed   bool   `json:"closed"`
}

type SSHSessionRequest struct {
	Command string `json:"command" validate:"required"`
}
