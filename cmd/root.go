/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"os"

	"github.com/openshift/ocm-container/pkg/engine"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	programName = "ocm-container"
)

var cfgFile string
var verbose bool
var containerEngine string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ocm-container",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running ocm-container")

		e := engine.New(containerEngine)

		out, err := e.Create("-it", "fedora:latest")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		err = e.StartAndAttach(out)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	configFileDefault := fmt.Sprintf("%s/%s/.%s.yaml", os.Getenv("HOME"), programName, programName)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", configFileDefault, "config file")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&containerEngine, "engine", "", "container engine to use (podman, docker)")

	rootCmd.MarkPersistentFlagRequired("engine")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".ocm-container" (without extension).
		viper.AddConfigPath(home + "/" + programName)
		viper.SetConfigType("yaml")
		viper.SetConfigName("." + programName)
		viper.AutomaticEnv()
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Fprintln(os.Stderr, "Error reading config file: %s", err)
	}
}
