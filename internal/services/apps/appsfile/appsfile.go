package appsfile

import (
	"encoding/json"
	"os"

	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/utils"
)

func New(file string) (*AppsService, error) {

	err := testFile(file)

	if err != nil {
		err = genFile()
	}

	if err != nil {
		return nil, err
	}

	return &AppsService{
		File: file,
	}, nil
}

func (as *AppsService) Table() (table models.AppsTable, err error) {
	as.mu.Lock()
	defer as.mu.Unlock()

	return as.read()
}

func (as *AppsService) Append(topic string, app models.AppDetails) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	table, err := as.read()
	if err != nil {
		return err
	}

	table[topic] = append(table[topic], app)

	return as.rewrite(table)
}

func (as *AppsService) Delete(topic string, topicQuery int) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	table, err := as.read()
	if err != nil {
		return err
	}

	if (len(table[topic])) < (topicQuery+1) || topicQuery < 0 {
		return ErrQueryOutOfRange
	}

	table[topic] = utils.RemoveSliceIndex(table[topic], topicQuery)

	if len(table[topic]) == 0 {
		delete(table, topic)
	}

	return as.rewrite(table)
}

// TODO: Append functional
func (as *AppsService) Edit(topic string, topicQuery int, app models.AppDetails) error {
	return nil
}

func (as *AppsService) read() (models.AppsTable, error) {

	table := make(models.AppsTable)

	data, err := os.ReadFile(as.File)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &table); err != nil {
		return nil, err
	}

	return table, nil
}

func (as *AppsService) rewrite(table models.AppsTable) error {

	file, err := os.OpenFile(as.File, os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(table, "", "   ")
	if err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}
