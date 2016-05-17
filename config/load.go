package config

import (
	"log"
	"path"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

// Load configuarion from a file
func Load(filename string) (*Config, error) {
	if filename != "" {
		log.Printf("[INFO] Loading configuration file: %s", filename)

		base := path.Base(filename)
		ext := path.Ext(filename)
		name := strings.TrimSuffix(base, ext)

		viper.SetConfigName(name)
		viper.AddConfigPath(path.Dir(filename))
	} else {
		viper.SetConfigFile(".watcher")
	}
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil && filename != "" {
		return nil, err
	}

	log.Printf("[INFO] Using config file: %s", viper.ConfigFileUsed())

	return SetDefaults()
}

// SetDefaults sets some defaults and parses some vaules
func SetDefaults() (cfg *Config, err error) {
	cfg = &Config{}

	viper.SetDefault("", Default)
	viper.Unmarshal(cfg)

	if cfg.Runtime.GOMAXPROCS == -1 {
		cfg.Runtime.GOMAXPROCS = runtime.NumCPU()
	}

	return cfg, nil
}
