package config

import "runtime"

// Default configuration
var Default = &Config{
	Server: Server{
		ListenAddr: ":9090",
	},
	Runtime: Runtime{
		GOGC:       800,
		GOMAXPROCS: runtime.NumCPU(),
	},
}
