package config

// Config specifies the application configuration
type Config struct {
	Server  Server
	Runtime Runtime
}

// Runtime Go configuration
type Runtime struct {
	GOGC       int
	GOMAXPROCS int
}

// AppConfig holds the current application configuration
var AppConfig Config
