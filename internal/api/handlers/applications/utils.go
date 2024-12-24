package applications

import (
	"encoding/json"
	"io"
	"os"

	"github.com/eterline/desky-backend/internal/utils"
)

func (as *AppsService) generateAppsFile() (AppsTable, error) {
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

func (as *AppsService) getAppTable() (table AppsTable, err error) {
	data, err := os.ReadFile(as.filePath)

	if err != nil {
		return as.generateAppsFile()
	}

	err = json.Unmarshal(data, &table)
	return table, err
}

func (as *AppsService) addApp(topic string, app AppDetails) error {
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

func (as *AppsService) rmApp(topic string, appNum int) error {
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

	table[topic] = utils.RemoveSliceIndex(table[topic], appNum)

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
