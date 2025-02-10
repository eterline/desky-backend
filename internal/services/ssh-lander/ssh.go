package sshlander

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

type AuthType int

const (
	_ AuthType = iota
	PrivateKeyMethod
	PasswordMethod
)

type SSHLanderService struct {
	config *ssh.ClientConfig

	stdinPipe  io.WriteCloser
	stdoutPipe io.Reader

	shellSession *ssh.Session
	shellСlient  *ssh.Client

	sync.Mutex
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

func (s *SSHLanderService) Connect(ctx context.Context, address string) error {
	cl, err := ssh.Dial("tcp", address, s.config)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	session, err := cl.NewSession()
	if err != nil {
		cl.Close()
		return fmt.Errorf("failed to create session: %w", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		cl.Close()
		return fmt.Errorf("request for pseudo terminal failed: %w", err)
	}

	inPipe, err := session.StdinPipe()
	if err != nil {
		session.Close()
		cl.Close()
		return fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	outPipe, err := session.StdoutPipe()
	if err != nil {
		inPipe.Close()
		session.Close()
		cl.Close()
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := session.Shell(); err != nil {
		inPipe.Close()
		session.Close()
		cl.Close()
		return fmt.Errorf("failed to start shell: %w", err)
	}

	s.stdinPipe = inPipe
	s.stdoutPipe = outPipe
	s.shellSession = session
	s.shellСlient = cl

	go func() {
		<-ctx.Done()
		s.Exit()
	}()

	return nil
}

func (s *SSHLanderService) SendCommand(command string) <-chan string {
	s.Lock()
	outputChan := make(chan string)

	go func() {
		defer close(outputChan)
		defer s.Unlock()

		endMarker := "COMMAND_FINISHED_MARKER"
		fullCommand := fmt.Sprintf("%s; pwd; echo %s\n", command, endMarker)

		_, err := fmt.Fprint(s.stdinPipe, fullCommand)
		if err != nil {
			outputChan <- fmt.Sprintf("stdin error: %v", err)
			return
		}

		reader := InitLineReader("", s.stdoutPipe)

		for {
			line := reader.Read()

			if strings.HasPrefix(line, endMarker) {
				return
			}

			if !strings.Contains(line, endMarker) {
				outputChan <- strings.TrimSpace(line)
			}
		}
	}()

	return outputChan
}

func (s *SSHLanderService) Exit() error {
	s.Lock()
	defer s.Unlock()

	if s.shellSession != nil {
		s.shellSession.Close()
	}
	if s.shellСlient != nil {
		s.shellСlient.Close()
	}
	return nil
}

type LineReader struct {
	filter *regexp.Regexp
	reader *bufio.Reader
}

func InitLineReader(regExp string, pipe io.Reader) *LineReader {

	filt := regexp.MustCompile(regExp)

	if regExp == "" {
		filt = regexp.MustCompile(`\\x1b\\[[0-9;]*[a-zA-Z]`)
	}

	return &LineReader{
		filter: filt,
		reader: bufio.NewReader(pipe),
	}
}

func (ls *LineReader) Read() string {
	line, err := ls.reader.ReadString('\n')
	if err != nil {
		if err != io.EOF {
			return fmt.Sprintf("stdout error: %v", err)
		}
	}

	filteredLine := ls.filter.ReplaceAllString(line, "")

	return strings.TrimSpace(filteredLine)
}
