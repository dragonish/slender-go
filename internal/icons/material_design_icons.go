package icons

import (
	"encoding/json"
	"os"
	"regexp"
	"slender/internal/logger"
	"slender/internal/model"
	"strings"
)

func materialDesignIconsBuilder() {
	MaterialDesignIcons = make(map[string]string)
	log := logger.New("file", model.ICONS_MATERIAL_DESIGN_ICONS_JS)

	fileRaw, err := os.ReadFile(model.ICONS_MATERIAL_DESIGN_ICONS_JS)
	if err == nil {
		var re = regexp.MustCompile(`(?m)var\s+mdi(\w+)\s*=\s*"(.+?)";`)
		for _, match := range re.FindAllStringSubmatch(string(fileRaw), -1) {
			MaterialDesignIcons[strings.ToLower(match[1])] = match[2]
		}
		log.Debug("build material-design-icons file to the map")
	} else {
		log.Err("unable to read icon file", err)
	}
}

func materialDesignIconsVersion() {
	MDIVer = "0.0.0"
	log := logger.New("file", model.ICONS_MATERIAL_DESIGN_ICONS_JSON)

	fileRaw, err := os.ReadFile(model.ICONS_MATERIAL_DESIGN_ICONS_JSON)
	if err == nil {
		var result map[string]any
		err = json.Unmarshal(fileRaw, &result)
		if err == nil {
			ver := result["version"].(string)
			if len(ver) > 0 {
				MDIVer = ver
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
