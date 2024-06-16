package icons

import "strings"

var (
	MaterialDesignIcons map[string]string
	SimpleIcons         map[string]string

	MDIVer string
	SIVer  string
)

// Build builds icon files.
func Build() {
	materialDesignIconsBuilder()
	simpleIconsBuilder()

	materialDesignIconsVersion()
	simpleIconsVersion()
}

// GetBuiltInIcon returns built-in icon path.
func GetBuiltInIcon(name string) string {
	if len(name) > 0 {
		if strings.HasPrefix(name, "mdi-") || strings.HasPrefix(name, "si-") {
			return "/assets/icons/" + name + ".svg"
		}
	}

	return ""
}
