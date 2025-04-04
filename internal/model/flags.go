package model

import (
	"slices"
	"strconv"
)

// Flags serves global state to the program.
type Flags struct {
	DebugMode   bool // enable debug mode.
	ShowVersion bool // show application version.
	ShowHelp    bool // show help document.

	PerformanceMode bool // enable performance mode.

	AccessPassword string // access password.
	AdminPassword  string // admin password (never empty). default: AccessPassword(not empty) or "p@$$w0rd"

	TokenAge uint16 // token max-age (days). minimum: 1

	LogLevel string // log output level.
	Port     uint16 // web service running port.

	Salt        string   // password salt.
	Secret      string   // JWT secret.
	AccessToken string   // access token.
	AdminToken  string   // admin token.
	LoginIDs    []string // login IDs.

	ServiceConfig string // service config file path.
}

// GetPortStr returns port string.
func (f *Flags) GetPortStr() string {
	return strconv.FormatUint(uint64(f.Port), 10)
}

// GetAccessCookieName returns access cookie name.
func (f *Flags) GetAccessCookieName() string {
	return COOKIE_ACCESS_PREFIX + f.GetPortStr()
}

// GetAdminCookieName returns admin cookie name.
func (f *Flags) GetAdminCookieName() string {
	return COOKIE_ADMIN_PREFIX + f.GetPortStr()
}

// GetTokenAgeSeconds returns token max-age seconds.
func (f *Flags) GetTokenAgeSeconds() int {
	return int(f.TokenAge) * 24 * 60 * 60
}

func (f *Flags) IsLogined(id string) bool {
	return slices.Contains(f.LoginIDs, id)
}
