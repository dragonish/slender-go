package icons

import (
	"encoding/json"
	"os"
	"regexp"
	"slender/internal/logger"
	"slender/internal/model"
	"strings"
)

func simpleIconsBuilder() {
	SimpleIcons = make(map[string]string)
	log := logger.New("file", model.ICONS_SIMPLE_ICONS_JS)

	fileRaw, err := os.ReadFile(model.ICONS_SIMPLE_ICONS_JS)
	if err == nil {
		var re = regexp.MustCompile(`[{,]si.+?:{title:.+?,slug:"(.+?)".+?},path:"(.+?)",`)
		for _, match := range re.FindAllStringSubmatch(string(fileRaw), -1) {
			SimpleIcons[strings.ToLower(match[1])] = match[2]
		}
		log.Debug("build simple-icons file to the map")
	} else {
		log.Err("unable to read icon file", err)
	}
}

func simpleIconsVersion() {
	SIVer = "0.0.0"
	log := logger.New("file", model.ICONS_SIMPLE_ICONS_JSON)

	fileRaw, err := os.ReadFile(model.ICONS_SIMPLE_ICONS_JSON)
	if err == nil {
		var result map[string]any
		err = json.Unmarshal(fileRaw, &result)
		if err == nil {
			ver := result["version"].(string)
			if len(ver) > 0 {
				SIVer = ver
			} else {
				log.Warn("unable to get version")
			}
		} else {
			log.Err("unable to parse info file", err)
		}
	} else {
		log.Err("unable to read info file", err)
	}
}
