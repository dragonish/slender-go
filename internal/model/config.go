package model

import "strings"

type ServiceConfig struct {
	AccessPassword *string `json:"accessPassword,omitempty" yaml:"access_password,omitempty"` // access password.
	AdminPassword  *string `json:"adminPassword,omitempty" yaml:"admin_password,omitempty"`   // admin password.

	LogLevel        *string `json:"logLevel,omitempty" yaml:"log_level,omitempty"`               // log output level.
	Port            *uint16 `json:"port,omitempty" yaml:"port,omitempty"`                        // web service running port.
	PerformanceMode *bool   `json:"performanceMode,omitempty" yaml:"performance_mode,omitempty"` // performance mode.

	TokenAge *uint16 `json:"tokenAge,omitempty" yaml:"token_age,omitempty"` // token max-age (days).
}

// UserConfig defines the user configuration field.
type UserConfig struct {
	Title        string `json:"title" yaml:"title"`                // website title.
	CustomFooter string `json:"customFooter" yaml:"custom_footer"` // custom footer.

	ShowSidebar     bool `json:"showSidebar" yaml:"show_sidebar"`          // show folders sidebar.
	ShowSearchInput bool `json:"showSearchInput" yaml:"show_search_input"` // show search input.
	ShowScrollTop   bool `json:"showScrollTop" yaml:"show_back_top"`       // show scroll to top button.

	ShowLatest  bool  `json:"showLatest" yaml:"show_latest"`   // show the module of latest added bookmarks.
	LatestTotal uint8 `json:"latestTotal" yaml:"latest_total"` // number of bookmarks in the latest module.
	ShowHot     bool  `json:"showHot" yaml:"show_hot"`         // show the module of hot bookmarks.
	HotTotal    uint8 `json:"hotTotal" yaml:"hot_total"`       // number of bookmarks in the hot module.

	UseLetterIcon   bool `json:"useLetterIcon" yaml:"use_letter_icon"`      // use first letter as icon.
	OpenInNewWindow bool `json:"openInNewWindow" yaml:"open_in_new_window"` // always open the bookmark in the new window.

	InternalNetwork string `json:"internalNetwork" yaml:"internal_network"` // Internal network address list. Multiple values are separated by commas (,).
}

// GetInternalNetwork returns the internal network address list as a slice of strings.
func (uc *UserConfig) GetInternalNetwork() []string {
	parts := strings.Split(uc.InternalNetwork, ",")

	var trimmedParts []string
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		if len(trimmedPart) > 0 {
			trimmedParts = append(trimmedParts, trimmedPart)
		}
	}

	return trimmedParts
}

// InOtherNetwork checks if the origin is in other networks.
func (uc *UserConfig) InOtherNetwork(origin string) bool {
	if len(uc.InternalNetwork) == 0 {
		return true
	}

	for _, network := range uc.GetInternalNetwork() {
		if strings.Contains(origin, network) {
			return false
		}
	}
	return true
}
