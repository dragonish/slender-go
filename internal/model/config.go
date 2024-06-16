package model

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
}
