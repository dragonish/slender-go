package config

import (
	"os"
	"slender/internal/data"
	"slender/internal/global"
	"slender/internal/logger"
	"slender/internal/model"

	"gopkg.in/yaml.v3"
)

// configFilename defines config filename.
var configFilename string

// DefaultConfigGenerator generates the default user configuration.
func DefaultConfigGenerator() model.UserConfig {
	return model.UserConfig{
		Title:           "Slender",
		CustomFooter:    `<div style="width: 100%; text-align: center;"><a href="https://github.com/dragonish/slender-go" target="_blank" rel="noopener">Slender</a> ©2024-2025 Created by <a href="https://github.com/dragonish" target="_blank" rel="noopener">dragonish</a></div>`,
		ShowSearchInput: true,
		ShowSidebar:     true,
		ShowScrollTop:   true,
		ShowLatest:      true,
		LatestTotal:     12,
		ShowHot:         true,
		HotTotal:        12,
		UseLetterIcon:   true,
		OpenInNewWindow: true,
	}
}

// readConfig loads config from readed config.
func readConfig(readedConfig *model.UserConfig) {
	global.Config = *readedConfig

	if global.Config.LatestTotal < 1 {
		global.Config.LatestTotal = 1
	} else if global.Config.LatestTotal > 100 {
		global.Config.LatestTotal = 100
	}

	if global.Config.HotTotal < 1 {
		global.Config.HotTotal = 1
	} else if global.Config.HotTotal > 100 {
		global.Config.HotTotal = 100
	}
}

// saveConfig saves config to file.
func saveConfig() {
	path := getConfigFilePath(configFilename)
	log := logger.New("path", path)

	res, err := yaml.Marshal(global.Config)
	if err == nil {
		wErr := os.WriteFile(path, res, os.ModePerm)
		if wErr != nil {
			log.Err("unable to write config to file", wErr)
		}
	} else {
		log.Err("unable to convert config content", err)
	}
}

// getConfigFilePath returns config file full path.
//
// The filename parameter does not take an extension name.
func getConfigFilePath(filename string) string {
	return model.DATA_DIR + "/" + filename + ".yaml"
}

// Load reads config from file.
//
// The filename parameter does not take an extension name.
func Load(filename string) {
	configFilename = filename

	path := getConfigFilePath(filename)
	log := logger.New("path", path)

	if data.IsPathExists(path) {
		content, err := os.ReadFile(path)
		if err == nil {
			var conf model.UserConfig
			pErr := yaml.Unmarshal(content, &conf)
			if pErr == nil {
				readConfig(&conf)
			} else {
				log.Err("unable to parse config file, default config is used", pErr)
				global.Config = DefaultConfigGenerator()
			}
		} else {
			log.Err("unable to read config file, default config is used", err)
			global.Config = DefaultConfigGenerator()
		}
	} else {
		log.Info("no config file found, default config is used")
		global.Config = DefaultConfigGenerator()
		//* write default config to file.
		saveConfig()
	}
}

// Update updates config.
func Update(conf model.UserConfig) {
	global.Config = conf
	saveConfig()
}

// PatchUpdate updates config.
func PatchUpdate(conf model.ConfigPatchBody) {
	if conf.Title != nil {
		global.Config.Title = *conf.Title
	}
	if conf.CustomFooter != nil {
		global.Config.CustomFooter = *conf.CustomFooter
	}
	if conf.ShowSidebar != nil {
		global.Config.ShowSidebar = *conf.ShowSidebar
	}
	if conf.ShowSearchInput != nil {
		global.Config.ShowSearchInput = *conf.ShowSearchInput
	}
	if conf.ShowScrollTop != nil {
		global.Config.ShowScrollTop = *conf.ShowScrollTop
	}
	if conf.ShowLatest != nil {
		global.Config.ShowLatest = *conf.ShowLatest
	}
	if conf.LatestTotal != nil {
		global.Config.LatestTotal = *conf.LatestTotal
	}
	if conf.ShowHot != nil {
		global.Config.ShowHot = *conf.ShowHot
	}
	if conf.HotTotal != nil {
		global.Config.HotTotal = *conf.HotTotal
	}
	if conf.UseLetterIcon != nil {
		global.Config.UseLetterIcon = *conf.UseLetterIcon
	}
	if conf.OpenInNewWindow != nil {
		global.Config.OpenInNewWindow = *conf.OpenInNewWindow
	}

	saveConfig()
}
