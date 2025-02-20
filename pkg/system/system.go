package system

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/bitfield/script"
)

func JqOutputCmd(cmdLine, parseParam string) *script.Pipe {
	if parseParam == "" {
		parseParam = "."
	}

	return script.Exec(cmdLine).Exec(fmt.Sprintf("jq %s", parseParam))
}

func ExecOut(pipe *script.Pipe) (out []byte, err error) {
	out, err = io.ReadAll(pipe)

	if err != nil {
		return nil, ErrExec(err)
	}

	if e := NewExitErrorCode(pipe.ExitStatus()); e != nil {
		return nil, e
	}

	return out, nil
}

func Uptime() UptimeDuration {

	out, err := io.ReadAll(script.Exec("uptime -p"))
	if err != nil {
		return 0
	}

	replacers := map[string]string{
		" hours":   "h",
		" days":    "d",
		" minutes": "m",
		", ":       "",
		"up ":      "",
		"\n":       "",
	}

	timeStr := string(out)

	for rep, dst := range replacers {
		timeStr = strings.ReplaceAll(timeStr, rep, dst)
	}

	t, _ := time.ParseDuration(timeStr)
	return UptimeDuration(t.Seconds())
}

func HostAddrs() AddrsList {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return AddrsList{"0.0.0.0"}
	}
	list := make(AddrsList, len(addrs))

	for i, addr := range addrs {
		list[i] = strings.Split(addr.String(), "/")[0]
	}

	return list
}

func ProcessCommand(cmd string) (output []byte) {
	output = []byte("err: can't display output")

	pipe := script.Exec(cmd)

	out, err := io.ReadAll(pipe)
	if err == nil {
		output = out
	}

	return output
}

func HandleCommand(data []byte) (*MessageCLI, error) {

	req := new(MessageCLI)

	if err := json.Unmarshal(data, &req); err != nil {
		return &MessageCLI{
			Command: "nil",
			Output:  "uncorrect request",
		}, ErrExec(err)
	}

	output := ProcessCommand(req.Command)

	return &MessageCLI{
		Command: req.Command,
		Output:  string(output),
	}, nil
}
