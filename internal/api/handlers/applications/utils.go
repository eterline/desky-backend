package applications

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/eterline/desky-backend/internal/utils"
)

func (as *AppsHandlerGroup) generateAppsFile() (AppsTable, error) {
	as.mutx.TryLock()
	defer as.mutx.Unlock()

	file, err := os.Create(as.filePath)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	data, err := json.MarshalIndent(ExampledtApps, "", "   ")
	if err != nil {
		return nil, err
	}

	_, err = file.Write(data)
	return ExampledtApps, err
}

func (as *AppsHandlerGroup) getAppTable() (table AppsTable, err error) {
	data, err := os.ReadFile(as.filePath)

	if err != nil {
		return as.generateAppsFile()
	}

	err = json.Unmarshal(data, &table)
	return table, err
}

func (as *AppsHandlerGroup) addApp(topic string, app AppDetails) error {
	as.mutx.TryLock()
	defer as.mutx.Unlock()

	if topic == "" {
		return handlers.NewErrorResponse(
			http.StatusBadRequest,
			errors.New("topic name can't be empty"),
		)
	}

	file, err := os.OpenFile(as.filePath, os.O_RDWR, 0644)
	defer file.Close()

	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	table := make(AppsTable)

	if err = json.Unmarshal(data, &table); err != nil {
		return err
	}

	table[topic] = append(table[topic], app)

	data, err = json.MarshalIndent(table, "", "   ")
	if err != nil {
		return err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = file.Write(data)

	return err
}

func (as *AppsHandlerGroup) rmApp(topic, appNum string) error {
	as.mutx.TryLock()
	defer as.mutx.Unlock()

	if topic == "" {
		return handlers.NewErrorResponse(
			http.StatusBadRequest,
			errors.New("topic name can't be empty"),
		)
	}

	appQueryNumber, err := strconv.Atoi(appNum)
	if err != nil {
		return handlers.NewErrorResponse(
			http.StatusBadRequest,
			errors.New("uncorrect app query number"),
		)
	}

	file, err := os.OpenFile(as.filePath, os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	table := make(AppsTable)

	if err = json.Unmarshal(data, &table); err != nil {
		return err
	}

	if (len(table[topic])) < (appQueryNumber+1) || appQueryNumber < 0 {
		return handlers.NewErrorResponse(
			http.StatusBadRequest,
			errors.New("app query number out of range"),
		)
	}

	table[topic] = utils.RemoveSliceIndex(table[topic], appQueryNumber)

	if len(table[topic]) == 0 {
		delete(table, topic)
	}

	data, err = json.MarshalIndent(table, "", "   ")
	if err != nil {
		return err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = file.Write(data)

	return err
}
