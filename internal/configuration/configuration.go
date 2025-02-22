package configuration

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
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
		return ErrUnknownExt
	}

	if err != nil {
		return ErrRead(err)
	}

	if err = config.Validation(); err != nil {
		return err
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

	file.Truncate(0)

	content, err := yaml.Marshal(DefaultParameters)
	if err != nil {
		return ErrMigration(err)
	}

	if _, err := file.Write(content); err != nil {
		return ErrMigration(err)
	}

	return nil
}

func (c *Configuration) Validation() error {
	s := validator.New()
	return s.Struct(c)
}
