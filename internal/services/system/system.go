package system

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
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

func HostAddrs() (AddrsList, error) {

	cmd, err := exec.LookPath("hostname")
	if err != nil {
		return nil, ErrExec(err)
	}

	pipe := script.Exec(fmt.Sprintf("%s -I", cmd))

	out, err := ExecOut(pipe)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(out), " "), nil
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

func HandleCommand(data []byte) (response *ResponseCLI, err error) {
	req := RequestCLI{}
	response = &ResponseCLI{}

	if err := json.Unmarshal(data, &req); err != nil {
		return &ResponseCLI{
			Command: "nil",
			Output:  "uncorrect request",
		}, ErrExec(err)
	}

	response.Command = req.Command
	response.Output = string(ProcessCommand(req.Command))

	return response, nil
}
