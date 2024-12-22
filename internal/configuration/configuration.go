package configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/stretchr/testify/assert/yaml"
)

var (
	instance *Configuration = nil
	once     sync.Once
)

func GetConfig() *Configuration {
	if instance == nil {
		return defaultConfig
	}
	return instance
}

var ConfigPathDefault = ""

// Initailazing app configuration.
// In error state returns default settings
func Init(path string) error {

	data, err := os.ReadFile(path)
	if err != nil {
		instance = defaultConfig
		return err
	}

	config := new(Configuration)

	switch filepath.Ext(path) {

	case "json":
		err = json.Unmarshal(data, config)
		break

	case "yml":
		err = yaml.Unmarshal(data, config)
		break

	case "yaml":
		err = yaml.Unmarshal(data, config)
		break

	default:
		instance = defaultConfig
		return ErrUnknownExt
	}

	if err != nil {
		instance = defaultConfig
		return err
	}

	instance = config
	return nil
}

// Server config functions =============================

// Return address string: IP:PORT => 0.0.0.0:3000
func (cfg *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%v", cfg.Connection.Addr, cfg.Connection.Port)
}

// Logger config function =============================
