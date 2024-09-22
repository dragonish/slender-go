package parse

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/pflag"

	"slender/internal/model"
	"slender/internal/version"
)

func parseCLI(baseFlags *model.Flags) model.Flags {
	var cliFlags = new(model.Flags)
	options := pflag.NewFlagSet("appFlags", pflag.ContinueOnError)
	options.SortFlags = false

	// Excute CLI command
	options.BoolVarP(&cliFlags.ShowVersion, model.KEY_VERSION, model.KEY_VERSION_SHORT, false, model.KEY_VERSION_DES)
	options.BoolVarP(&cliFlags.ShowHelp, model.KEY_HELP, model.KEY_HELP_SHORT, false, model.KEY_HELP_DES)

	options.BoolVarP(&cliFlags.DebugMode, model.KEY_DEBUG, model.KEY_DEBUG_SHORT, false, model.KEY_DEBUG_DES)

	options.BoolVarP(&cliFlags.PerformanceMode, model.KEY_PERFORMANCE_MODE, model.KEY_PERFORMANCE_MODE_SHORT, baseFlags.PerformanceMode, model.KEY_PERFORMANCE_MODE_DES)

	options.StringVarP(&cliFlags.AccessPassword, model.KEY_ACCESS_PWD, model.KEY_ACCESS_PWD_SHORT, baseFlags.AccessPassword, model.KEY_ACCESS_PWD_DES)
	options.StringVarP(&cliFlags.AdminPassword, model.KEY_ADMIN_PWD, model.KEY_ADMIN_PWD_SHORT, baseFlags.AdminPassword, model.KEY_ADMIN_PWD_DES)

	options.Uint16VarP(&cliFlags.TokenAge, model.KEY_TOKEN_AGE, model.KEY_TOKEN_AGE_SHORT, baseFlags.TokenAge, model.KEY_TOKEN_AGE_DES)

	options.StringVarP(&cliFlags.LogLevel, model.KEY_LOG_LEVEL, model.KEY_LOG_LEVEL_SHORT, baseFlags.LogLevel, model.KEY_LOG_LEVEL_DES)
	options.Uint16VarP(&cliFlags.Port, model.KEY_PORT, model.KEY_PORT_SHORT, baseFlags.Port, model.KEY_PORT_DES)

	options.StringVarP(&cliFlags.ServiceConfig, model.KEY_CONFIG_FILE, model.KEY_CONFIG_FILE_SHORT, baseFlags.ServiceConfig, model.KEY_CONFIG_FILE_DES)

	_ = options.Parse(os.Args)

	if excuteCLI(cliFlags, options) {
		os.Exit(0)
	}

	//* Parse service config file.
	if cliFlags.ServiceConfig != "" {
		lErr := loadServiceConfig(baseFlags, cliFlags.ServiceConfig)
		if lErr == nil {
			if !isManual(options, model.KEY_ACCESS_PWD) {
				cliFlags.AccessPassword = baseFlags.AccessPassword
			}
			if !isManual(options, model.KEY_ADMIN_PWD) {
				cliFlags.AdminPassword = baseFlags.AdminPassword
			}
			if !isManual(options, model.KEY_LOG_LEVEL) {
				cliFlags.LogLevel = baseFlags.LogLevel
			}
			if !isManual(options, model.KEY_PORT) {
				cliFlags.Port = baseFlags.Port
			}
			if !isManual(options, model.KEY_PERFORMANCE_MODE) {
				cliFlags.PerformanceMode = baseFlags.PerformanceMode
			}
			if !isManual(options, model.KEY_TOKEN_AGE) {
				cliFlags.TokenAge = baseFlags.TokenAge
			}
		} else {
			//? Empty the property to identify that no service config is applied.
			cliFlags.ServiceConfig = ""
		}
	}

	return *cliFlags
}

func excuteCLI(cliFlags *model.Flags, options *pflag.FlagSet) bool {
	verInfo := getVersion()

	if cliFlags.ShowVersion {
		fmt.Println(verInfo)
		return true
	}

	if cliFlags.ShowHelp {
		fmt.Println(verInfo)
		fmt.Println("Commands:")
		options.PrintDefaults()
		return true
	}

	return false
}

func getVersion() (info string) {
	info = fmt.Sprintf("Slender v%s-%s %s/%s BuildDate=%s", version.Version, strings.ToLower(version.Commit), runtime.GOOS, runtime.GOARCH, version.BuildDate)

	return
}

func isManual(pf *pflag.FlagSet, key string) bool {
	flag := pf.Lookup(key)
	return flag.Changed
}
