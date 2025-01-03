package configuration

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-playground/validator/v10"
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

	switch filepath.Ext(path)[1:] {

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

	err = validator.New().Struct(config)

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

func (cfg *ServerConfig) JWTSecretBytes() []byte {
	return []byte(cfg.JWTSecret)
}

func (cfg *ServerConfig) PageAddr() string {

	u := url.URL{

		Scheme: func(v bool) string {
			if v {
				return "https"
			}
			return "http"
		}(cfg.TLS.Enabled),

		Host: cfg.Connection.Hostname + func(port uint16) string {
			if port == 80 || port == 443 {
				return ""
			}
			return fmt.Sprintf(":%v", port)
		}(cfg.Connection.Port),
		Path: "/",
	}

	return u.String()
}

// Logger config function =============================
