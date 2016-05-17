// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/harnash/watcher/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch server",
	Long:  `Launch server with specified config.`,
	Run: func(cmd *cobra.Command, args []string) {
		var conf server.Config

		viper.Unmarshal(&conf)
		server.Run(conf)
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
	RootCmd.PersistentFlags().IntP("port", "p", 8080, "Port on which server should listen for connections")
	RootCmd.PersistentFlags().StringP("address", "a", "", "Address on which server should listen for connections")

	cobra.OnInitialize(initServerConfig)
}

func initServerConfig() {
	server.LoadDefaults()

	viper.BindPFlag("server.listenport", RootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("server.listenaddress", RootCmd.PersistentFlags().Lookup("address"))
}
