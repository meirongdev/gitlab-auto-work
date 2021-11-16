/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	"github.com/chengshidaomin/gitlab-auto-work/internal"
	"github.com/chengshidaomin/gitlab-auto-work/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

var workConfig config.WorkConfig

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Info().Msg("Start to do all things from config file")

		if err := viper.Unmarshal(&workConfig); err != nil {
			log.Error().Msgf("Config file[%s] is invalid", viper.ConfigFileUsed())
		}
		log.Debug().Msgf("config info: %+v", workConfig)

		git, err := gitlab.NewClient(workConfig.Token, gitlab.WithBaseURL(workConfig.BaseUrl))
		if err != nil {
			log.Error().Msgf("Failed to create client: %v", err)
			os.Exit(2)
		}
		var workflow = &internal.Workflow{
			WorkConfig: workConfig,
			Client:     git,
			Log:        &log.Logger,
		}

		workflow.Run()
	},
}

func init() {
	rootCmd.AddCommand(allCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// allCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// allCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
