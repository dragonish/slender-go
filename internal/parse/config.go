package parse

import (
	"encoding/json"
	"os"
	"slender/internal/data"
	"slender/internal/logger"
	"slender/internal/model"
	"strings"

	"gopkg.in/yaml.v3"
)

// loadServiceConfig loads the service configuration file.
func loadServiceConfig(baseFlags *model.Flags, path string) error {
	log := logger.New("path", path)

	if data.IsPathExists(path) {
		content, err := os.ReadFile(path)
		if err == nil {
			conf, rErr := parseServiceConfig(&content, strings.HasSuffix(path, ".json"))
			if rErr == nil {
				serviceConfigHandler(baseFlags, &conf)
			} else {
				return log.Err("unable to parse service config file", rErr)
			}
		} else {
			return log.Err("unable to read service config file", err)
		}
	} else {
		return log.NewErr("no service config file found")
	}

	return nil
}

// parseServiceConfig parses service config from file content.
func parseServiceConfig(content *[]byte, isJSONFile bool) (conf model.ServiceConfig, err error) {
	if isJSONFile {
		err = json.Unmarshal(*content, &conf)
	} else {
		// Default as YAML file
		err = yaml.Unmarshal(*content, &conf)
	}

	return
}

// serviceConfigHandler handles the assignment of the config.
func serviceConfigHandler(baseFlags *model.Flags, conf *model.ServiceConfig) {
	if conf.AccessPassword != nil {
		baseFlags.AccessPassword = *conf.AccessPassword
	}
	if conf.AdminPassword != nil {
		baseFlags.AdminPassword = *conf.AdminPassword
	}
	if conf.LogLevel != nil {
		baseFlags.LogLevel = *conf.LogLevel
	}
	if conf.Port != nil {
		baseFlags.Port = *conf.Port
	}
	if conf.PerformanceMode != nil {
		baseFlags.PerformanceMode = *conf.PerformanceMode
	}
	if conf.TokenAge != nil {
		baseFlags.TokenAge = *conf.TokenAge
	}
}
