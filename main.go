package main

import (
	"github.com/harnash/watcher/cmd"
	"github.com/spf13/viper"
)

var (
	Version   string
	BuildTime string
)

func main() {
	viper.Set("version", Version)
	viper.Set("buildtime", BuildTime)

	cmd.Execute()
}
