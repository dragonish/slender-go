package model

// Fields and descriptions of the flag.
const (
	KEY_DEBUG       = "debug"
	KEY_DEBUG_SHORT = "D"
	KEY_DEBUG_DES   = "Enable debug mode"

	KEY_VERSION       = "version"
	KEY_VERSION_SHORT = "v"
	KEY_VERSION_DES   = "Show application version"

	KEY_HELP       = "help"
	KEY_HELP_SHORT = "h"
	KEY_HELP_DES   = "Show help document"

	KEY_ACCESS_PWD       = "access_pwd"
	KEY_ACCESS_PWD_SHORT = "a"
	KEY_ACCESS_PWD_DES   = "Specify access password"

	KEY_ADMIN_PWD       = "admin_pwd"
	KEY_ADMIN_PWD_SHORT = "d"
	KEY_ADMIN_PWD_DES   = "Specify admin password"

	KEY_TOKEN_AGE       = "token_age"
	KEY_TOKEN_AGE_SHORT = "t"
	KEY_TOKEN_AGE_DES   = "Specify token age (days)"

	KEY_LOG_LEVEL       = "log"
	KEY_LOG_LEVEL_SHORT = "l"
	KEY_LOG_LEVEL_DES   = "Specify log output level. Optional: Debug|Info|Warn|Error"

	KEY_PORT       = "port"
	KEY_PORT_SHORT = "p"
	KEY_PORT_DES   = "Specify web service running port"
)
