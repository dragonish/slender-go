package parse

import (
	"fmt"
	"slender/internal/data"
	"slender/internal/global"
	"slender/internal/logger"
)

// Parse handles environment variables and start commands.
//
// Assign values to the global flag state.
func Parse() {
	envs := parseEnvVars()
	flags := parseCLI(&envs)

	// Configure admin password.
	if flags.AdminPassword == "" {
		//? If admin password is unset, default to access password(not empty) or "p@$$w0rd".
		if flags.AccessPassword != "" {
			flags.AdminPassword = flags.AccessPassword
		} else {
			flags.AdminPassword = "p@$$w0rd"
		}
	}

	// Configure token age.
	if flags.TokenAge < 1 {
		flags.TokenAge = 1
	}

	fmt.Println(getVersion())

	if flags.DebugMode {
		flags.LogLevel = "Debug"
		logger.Info("enable debug mode", "debug_mode", flags.DebugMode)
	}

	if flags.PerformanceMode {
		logger.Info("enable performance mode", "performance_mode", flags.PerformanceMode)
	}

	logger.Info("access password", "access_pwd", data.MaskWithStars(flags.AccessPassword))
	logger.Info("admin password", "admin_pwd", data.MaskWithStars(flags.AdminPassword))
	logger.Info("token age", "days", flags.TokenAge)
	logger.Info("log output level", "log_level", flags.LogLevel)
	logger.Info("web service running port", "port", flags.Port)

	//? Set the lowest log output level after output basic information.
	logger.SetLevel(flags.LogLevel)

	global.Flags = flags
}
