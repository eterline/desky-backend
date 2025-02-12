package sshlander

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

type SessionCredentials interface {
	ValueUser() string
	Socket() string
	UsePrivateKey() bool
	Password() string
	PrivateKey() []byte
}

func newConfig(creds SessionCredentials, HostKeyCallBack ssh.HostKeyCallback) *ssh.ClientConfig {

	var sshAuthMethods []ssh.AuthMethod

	signer, err := ssh.ParsePrivateKey(creds.PrivateKey())
	if creds.UsePrivateKey() && err == nil {
		sshAuthMethods = append(sshAuthMethods, ssh.PublicKeys(signer))
	}

	sshAuthMethods = append(sshAuthMethods, ssh.Password(creds.Password()))

	return &ssh.ClientConfig{
		User:            creds.ValueUser(),
		Auth:            sshAuthMethods,
		HostKeyCallback: HostKeyCallBack,
	}
}

func NewClientSession(creds SessionCredentials, hostCallBack ssh.HostKeyCallback, uuid uuid.UUID) (*SSHSession, error) {

	tcpDial, err := ssh.Dial("tcp", creds.Socket(), newConfig(creds, hostCallBack))
	if err != nil {
		return nil, NewError(uuid, fmt.Sprintf("tcp dial error: %v", err))
	}

	sshSession, err := tcpDial.NewSession()
	if err != nil {
		tcpDial.Close()
		return nil, NewError(uuid, fmt.Sprintf("failed to create session: %v", err))
	}

	return &SSHSession{
		credentials: creds,
		uuid:        uuid,

		sshClient:  tcpDial,
		sshSession: sshSession,
	}, nil
}
