package sshlander

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

var (
	LineBreak      = []byte{0x0A}
	CarriageReturn = []byte{0x0D}
	ExitLine       = []byte{65, 78, 69, 74}
)

const (
	CommandDetermit = "DESKY_DETERMIT_COMMANDLINE"
	ExitCommand     = "exit"
)

type TerminalType string

const (
	XtermColored TerminalType = "xterm-256color"
	Xterm        TerminalType = "xterm"
)

type TerminalSession struct {
	ctx    context.Context
	closer context.CancelFunc

	stdIn  io.WriteCloser
	stdOut io.Reader
	stdErr io.Reader

	*SSHSession
	sync.Mutex
}

func ConnectTerminal(ctx context.Context, session *SSHSession, terminal TerminalType) (*TerminalSession, error) {

	terminalSettings := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.sshSession.RequestPty(string(terminal), 80, 40, terminalSettings); err != nil {
		fmt.Printf("err ssh: %v\n", err)
		session.CloseDial()
		return nil, NewError(session.uuid, fmt.Sprintf("xterm connection error: %v", err))
	}

	inPipe, err := session.sshSession.StdinPipe()
	if err != nil {
		fmt.Printf("err ssh: %v\n", err)
		session.CloseDial()
		return nil, NewError(session.uuid, fmt.Sprintf("failed to get stdin pipe: %v", err))
	}

	outPipe, err := session.sshSession.StdoutPipe()
	if err != nil {
		fmt.Printf("err ssh: %v\n", err)
		inPipe.Close()
		session.CloseDial()
		return nil, NewError(session.uuid, fmt.Sprintf("failed to get stdout pipe: %v", err))
	}

	errPipe, err := session.sshSession.StderrPipe()
	if err != nil {
		fmt.Printf("err ssh: %v\n", err)
		inPipe.Close()
		session.CloseDial()
		return nil, NewError(session.uuid, fmt.Sprintf("failed to get stderr pipe: %v", err))
	}

	if err := session.sshSession.Shell(); err != nil {
		fmt.Printf("err ssh: %v\n", err)
		inPipe.Close()
		session.CloseDial()
		return nil, NewError(session.uuid, fmt.Sprintf("failed to start session shell: %v", err))
	}

	ctx, closer := context.WithCancel(ctx)

	termSession := &TerminalSession{
		ctx:    ctx,
		closer: closer,

		stdIn:  inPipe,
		stdOut: outPipe,
		stdErr: errPipe,

		SSHSession: session,
	}

	return termSession, nil
}

func (term *TerminalSession) TerminalErrRead() <-chan []byte {
	readerChan := make(chan []byte, 1)

	go func() {
		defer close(readerChan)

		r := bufio.NewReader(term.stdErr)

		for {
			select {
			case <-term.ctx.Done():
				return

			default:

				line, err := r.ReadString('\n')
				if err != nil {
					readerChan <- []byte(fmt.Sprintf("stderr error: %v", err))
					return
				}

				readerChan <- []byte(line)
			}
		}
	}()

	return readerChan
}

func (term *TerminalSession) FromTerminalLines(w io.Writer, lineDetermiter byte) error {

	r := bufio.NewReader(term.stdOut)

	for {
		select {
		case <-term.ctx.Done():
			return nil

		default:

			line, err := r.ReadString(lineDetermiter)
			if err != nil {
				return err
			}

			if _, err := w.Write([]byte(line)); err != nil {
				return err
			}
		}
	}
}

func (term *TerminalSession) FromTerminalBytes(w io.Writer, wrSize int64) error {
	for {
		select {

		case <-term.ctx.Done():
			return nil

		default:

			if _, err := io.CopyN(w, term.stdOut, wrSize); err != nil {
				return err
			}
		}
	}
}

func (term *TerminalSession) TerminalRead() <-chan []byte {
	readerChan := make(chan []byte, 10)

	go func() {
		defer close(readerChan)

		r := bufio.NewReader(term.stdOut)

		for {
			select {
			case <-term.ctx.Done():
				return

			default:

				line, err := r.ReadString('\n')
				if err != nil {
					readerChan <- []byte(fmt.Sprintf("stdout error: %v", err))
					return
				}

				if !strings.Contains(line, CommandDetermit) {
					readerChan <- []byte(line)
				}
			}
		}
	}()

	return readerChan
}

// =============== Terminal writing =======================

// WriteLineBreak - sends '\n' to stdin
func (term *TerminalSession) WriteLineBreak() error {
	if _, err := term.stdIn.Write(LineBreak); err != nil {
		return NewError(term.uuid, fmt.Sprintf("failed to enter command: %v", err))
	}
	return nil
}

// WriteLineBreak - sends 'exit' to stdin for session exit
func (term *TerminalSession) WriteExit() error {
	if _, err := term.stdIn.Write(ExitLine); err != nil {
		return NewError(term.uuid, fmt.Sprintf("failed to send 'exit': %v", err))
	}
	return nil
}

// SendCleared - send stream to terminal stdin
func (term *TerminalSession) Send(data []byte) error {

	term.Lock()
	defer term.Unlock()

	if term.ctx.Err() != nil {
		return NewError(term.uuid, fmt.Sprintf("terminal session is closed"))
	}

	if _, err := term.stdIn.Write(data); err != nil {
		return NewError(term.uuid, fmt.Sprintf("failed to send string: %v", err))
	}

	if err := term.WriteLineBreak(); err != nil {
		return err
	}

	return nil
}

// SendCleared - send stream to terminal stdin with something end
// Example: '; pwd; echo %s", data, CommandDetermit'
func (term *TerminalSession) SendWithPostfix(data []byte, postfix string) error {

	term.Lock()
	defer term.Unlock()

	if term.ctx.Err() != nil {
		return NewError(term.uuid, fmt.Sprintf("terminal session is closed"))
	}

	if _, err := term.stdIn.Write(append(data, []byte(postfix)...)); err != nil {
		return NewError(term.uuid, fmt.Sprintf("failed to send string: %v", err))
	}

	if err := term.WriteLineBreak(); err != nil {
		return err
	}

	return nil
}

func (term *TerminalSession) Exit() error {

	term.stdIn.Close()

	term.closer()

	return term.CloseDial()
}
