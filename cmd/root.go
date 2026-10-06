package cmd

import (
	"fmt"
	"log"
	"os"

	"icmongolang/config"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var (
	cfgFile string
	RootCmd = &cobra.Command{
		Use:   "icmongolang",
		Short: "IoT Monitoring System",
		Long:  "IoT Monitoring System with MQTT, InfluxDB, WebSocket",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
)

// Execute runs the root command (called from package main).
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-base.yaml)")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.AddCommand(serveCmd)
	RootCmd.AddCommand(migrateCmd)
	RootCmd.AddCommand(workerCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfgViper, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	_, err = config.ParseConfig(cfgViper)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}
}
