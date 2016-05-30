package server

// Config holds the server configuration
type Config struct {
	// listenAddress specifies address on which server should listen for new connnections with port
	ListenAddress string `mapstructure:"server.ListenAddress"`
}
