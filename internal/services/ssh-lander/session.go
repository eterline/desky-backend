package sshlander

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

type SSHSession struct {
	credentials SessionCredentials
	uuid        uuid.UUID

	sshSession *ssh.Session
	sshClient  *ssh.Client
}

func (ss *SSHSession) UUID() string {
	return ss.uuid.String()
}

func (ss *SSHSession) Instance() string {
	return fmt.Sprintf(
		"%s@%s", ss.credentials.ValueUser(), ss.credentials.Socket(),
	)
}

func (ss *SSHSession) CloseDial() (err error) {
	if ss.sshClient != nil {
		err = ss.sshClient.Close()
	}
	if ss.sshSession != nil {
		err = ss.sshSession.Close()
	}
	return err
}
