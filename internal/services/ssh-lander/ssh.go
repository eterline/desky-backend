package sshlander

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/ssh"
)

type AuthType int

const (
	_ AuthType = iota
	PrivateKeyMethod
	PasswordMethod
)

type SSHLanderService struct {
	config       *ssh.ClientConfig
	shellSession *ssh.Session

	stdinPipe  io.WriteCloser
	stdoutPipe io.Reader
}

func New(user string) *SSHLanderService {

	return &SSHLanderService{
		config: &ssh.ClientConfig{
			User: user,
			Auth: make([]ssh.AuthMethod, 0),
		},
		shellSession: nil,
	}
}

func (s *SSHLanderService) SetupAuth(t AuthType, secret string) error {

	s.config.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	switch t {

	case PasswordMethod:
		s.config.Auth = append(s.config.Auth, ssh.Password(secret))
		return nil

	case PrivateKeyMethod:
		key, err := ssh.ParsePrivateKey([]byte(secret))
		if err != nil {
			return err
		}

		s.config.Auth = append(s.config.Auth, ssh.PublicKeys(key))
		return nil

	default:
		return errors.New("unknown auth method")
	}

}

func (s *SSHLanderService) Connect(address string) error {
	cl, err := ssh.Dial("tcp", address, s.config)
	if err != nil {
		return err
	}

	session, err := cl.NewSession()
	if err != nil {
		return err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // Включение echo
		ssh.TTY_OP_ISPEED: 14400, // Скорость передачи
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return fmt.Errorf("request for pseudo terminal failed: %v", err)
	}

	stdinPipe, err := session.StdinPipe()
	if err != nil {
		return err
	}

	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	if err := session.Shell(); err != nil {
		return err
	}

	s.stdinPipe = stdinPipe
	s.stdoutPipe = stdoutPipe

	s.shellSession = session
	return nil
}

// TODO: доделать работу с ssh
func (s *SSHLanderService) SendCommand(command string) (string, error) {
	endMarker := "COMMAND_FINISHED_MARKER_12345"
	fullCommand := fmt.Sprintf("%s; echo %s\n", command, endMarker)

	_, err := fmt.Fprint(s.stdinPipe, fullCommand)
	if err != nil {
		return "", err
	}

	var outputBuf bytes.Buffer
	reader := bufio.NewReader(s.stdoutPipe)

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		outputBuf.WriteString(line)

		if strings.Contains(line, endMarker) {
			break
		}
	}

	finalOutput := cleanOutput(outputBuf.String(), endMarker)
	return finalOutput, nil
}

func cleanOutput(rawOutput, marker string) string {
	lines := strings.Split(rawOutput, "\n")
	var result []string

	for _, line := range lines {
		if strings.Contains(line, "Welcome to") ||
			strings.Contains(line, "Documentation:") ||
			strings.Contains(line, marker) ||
			strings.HasPrefix(line, "root@") ||
			strings.TrimSpace(line) == "" {
			continue
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func (s *SSHLanderService) Exit() error {
	return s.shellSession.Close()
}
