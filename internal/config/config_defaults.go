package config

import (
	"runtime"
)

var configDefaults = &Configuration{
	Runtime: RuntimeConfig{
		GOGC:       100,
		GOMAXPROCS: runtime.NumCPU(),
	},
	Listen: Listen{
		Host: "127.0.0.1",
		Port: 8080,
	},
}
