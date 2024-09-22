package parse

import (
	"fmt"

	"github.com/caarlos0/env/v10"

	"slender/internal/model"
)

func parseEnvVars() (stor model.Flags) {
	cfg := model.Env{}

	opts := env.Options{
		Prefix: model.ENV_VAR_PREFIX,
	}

	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		fmt.Printf("Error parsing environment variables: %s", err.Error())
		return
	}

	stor.AccessPassword = cfg.AccessPassword
	stor.AdminPassword = cfg.AdminPassword
	stor.TokenAge = cfg.TokenAge
	stor.LogLevel = cfg.LogLevel
	stor.Port = cfg.Port
	stor.PerformanceMode = cfg.PerformanceMode == 1
	stor.ServiceConfig = cfg.ServiceConfig

	return
}
