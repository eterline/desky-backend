package configuration

import (
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var instance *Configuration

func GetConfig() *Configuration {

	if instance == nil {
		return DefaultParameters
	}
	return instance
}

func Init(path string) error {

	file, err := os.ReadFile(path)
	if err != nil {
		return ErrRead(err)
	}

	config := &Configuration{}

	switch filepath.Ext(path)[1:] {

	case "yml":
		err = yaml.Unmarshal(file, config)
		break

	case "yaml":
		err = yaml.Unmarshal(file, config)
		break

	default:
		if err != nil {
			return ErrUnknownExt
		}
	}

	if err != nil {
		return ErrRead(err)
	}

	instance = config

	return nil
}

func Migrate(fileName string, filePermit fs.FileMode) error {

	file, err := os.OpenFile(
		fileName,
		os.O_CREATE|os.O_RDWR,
		filePermit,
	)

	defer file.Close()

	if err != nil {
		return ErrMigration(err)
	}

	content, err := yaml.Marshal(DefaultParameters)
	if err != nil {
		return ErrMigration(err)
	}

	if _, err := file.Write(content); err != nil {
		return ErrMigration(err)
	}

	return nil
}

func (c *Configuration) Valid() error {
	return nil
}
