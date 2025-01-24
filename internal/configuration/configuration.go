package configuration

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	instance *Configuration = nil
)

func GetConfig() *Configuration {
	return instance
}

func MustConfig(path string) {

	file, err := os.ReadFile(path)
	if err != nil {
		panic(ErrRead(err))
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
		panic(ErrUnknownExt)
	}

	if err != nil {
		panic(ErrRead(err))
	}

	instance = config
}

func GenerateFile(name string) error {

	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0644)
	defer file.Close()

	if err != nil {
		return err
	}

	content, err := yaml.Marshal(DefaultParameters)
	if err != nil {
		return err
	}

	if _, err := file.Write(content); err != nil {
		return err
	}

	return nil
}
