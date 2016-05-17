package server

import "github.com/spf13/viper"

// Config holds the server configuration
type Config struct {
	// ListenPort speicifies the port server will bind to
	ListenPort int `mapstructure:"server.ListenPort"`
	// listenAddress specifies address on which server should listen for new connnections
	ListenAddress string `mapstructure:"server.ListenAddress"`
}

// LoadDefaults will set the default for the server configuration
func LoadDefaults() {
	viper.SetDefault("server.ListenAddress", "")
	viper.SetDefault("server.ListenPort", 8080)
}
