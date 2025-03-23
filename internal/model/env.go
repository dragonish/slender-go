package model

// Env storages environment variable.
type Env struct {
	AccessPassword string `env:"ACCESS_PWD,unset" envDefault:""` // access password. default: ""
	AdminPassword  string `env:"ADMIN_PWD,unset" envDefault:""`  // admin password. default: ""

	LogLevel        string `env:"LOG_LEVEL" envDefault:"Info"`     // log output level. default: "Info"
	Port            uint16 `env:"PORT" envDefault:"8080"`          // web service running port. default: 8080
	PerformanceMode string `env:"PERFORMANCE_MODE" envDefault:"0"` // performance mode. default: "0"

	TokenAge uint16 `env:"TOKEN_AGE" envDefault:"30"` // token age (days). default: 30

	ServiceConfig string `env:"SERVICE_CONFIG" envDefault:""` // service config file path. default: ""
}
