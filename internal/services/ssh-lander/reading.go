package sshlander

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

const (
	FilterASCII = `\\x1b\\[[0-9;]*[a-zA-Z]`
)

type PipeWorker struct {
	determit string
	filter   *regexp.Regexp

	reader *bufio.Reader
	writer io.Writer
}

func NewPipeWorker(
	filteringExpression, endDeterminant string,
	r io.Reader, w io.Writer,
) *PipeWorker {

	filter := regexp.MustCompile(filteringExpression)

	return &PipeWorker{
		determit: fmt.Sprintf("echo %s", endDeterminant),
		filter:   filter,

		reader: bufio.NewReader(r),
		writer: bufio.NewWriter(w),
	}
}

func (ls *PipeWorker) ReadLine(delim byte) (outputLine string) {
	line, err := ls.reader.ReadString(delim)
	if err != nil {
		return fmt.Sprintf("stdout error: %v", err)
	}

	fmt.Println(line)

	outputLine = ls.filter.ReplaceAllString(line, "")
	// outputLine = strings.TrimSpace(outputLine)

	return
}

// func (ls *PipeWorker) ReadLineWithout(
// 	delim byte,
// 	stringTest func(string, string) bool,
// 	substr string,
// ) (string, bool) {
// 	readLine := ls.ReadLine(delim)
// 	return readLine, !stringTest(readLine, substr)
// }

func joinInput(end string, values ...string) string {
	return strings.Join(values, "; ")
}

func (ls *PipeWorker) WriteStdin(value string) error {
	_, err := fmt.Fprintf(ls.writer, "%s\n", value)

	if err != nil {
		return fmt.Errorf("stdin error: %v", err)
	}

	return nil
}
