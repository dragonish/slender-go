package model

import "strconv"

// Flags serves global state to the program.
type Flags struct {
	DebugMode   bool // enable debug mode.
	ShowVersion bool // show application version.
	ShowHelp    bool // show help document.

	PerformanceMode bool // enable performance mode.

	AccessPassword string // access password.
	AdminPassword  string // admin password. default: AccessPassword(not empty) or "p@$$w0rd"

	TokenAge uint16 // token age (days). minimum: 1

	LogLevel string // log output level.
	Port     uint16 // web service running port.

	Salt        string // password salt.
	Secret      string // JWT secret.
	AccessToken string // access token.
	AdminToken  string // admin token.
}

// GetPortStr returns port string.
func (f *Flags) GetPortStr() string {
	return strconv.FormatUint(uint64(f.Port), 10)
}

// GetTokenAgeSeconds returns token age seconds.
func (f *Flags) GetTokenAgeSeconds() int {
	return int(f.TokenAge) * 24 * 60 * 60
}
