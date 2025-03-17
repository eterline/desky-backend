package controllers

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/eterline/desky-backend/internal/services/handler"
)

type LoggingCollector interface {
	Clean()
	GetStack() ([]string, error)
}

type ParametersControllers struct {
	logCollect LoggingCollector
}

func InitParameters(collect LoggingCollector) *ParametersControllers {
	return &ParametersControllers{
		logCollect: collect,
	}
}

func (pc *ParametersControllers) GetLogs(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "parameters.get-logs"
	r.ParseForm()

	logsRead, err := scanFile("./logs/trace.log")
	if err != nil {
		return op, err
	}

	linesNumber, _ := strconv.Atoi(r.FormValue("lines"))

	if linesNumber > 1 && linesNumber-1 <= len(logsRead) {
		logsRead = logsRead[(len(logsRead) - linesNumber):]
	}

	return op, handler.WriteJSON(w, http.StatusOK, logsRead)
}

func (pc *ParametersControllers) Errors(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "parameters.errors"

	stack, err := pc.logCollect.GetStack()
	if err != nil {
		return op, err
	}

	return op, handler.WriteJSON(w, http.StatusOK, stack)
}

func scanFile(logfile string) ([]string, error) {

	var sliceString []string

	f, err := os.OpenFile(logfile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, handler.NewErrorResponse(
			http.StatusNotImplemented,
			fmt.Errorf("can not open logs file. err: %v", err),
		)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		sliceString = append(sliceString, sc.Text())
	}
	if err := sc.Err(); err != nil {
		return nil, handler.NewErrorResponse(
			http.StatusNotImplemented,
			fmt.Errorf("can not scan logs file. err: %v", err),
		)
	}

	return sliceString, nil
}
