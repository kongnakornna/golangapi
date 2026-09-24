package cmd

import (
	"fmt"
	"log"
	"os"

	"icmongolang/config"

	"github.com/spf13/cobra"
)

var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "Run Kafka-related commands",
}

// RootCmd represents the base command when called without any subcommands
var (
	cfgFile string
	RootCmd = &cobra.Command{
		Use:   "icmongolang",
		Short: "IoT Monitoring System",
		Long:  "IoT Monitoring System with MQTT, InfluxDB, WebSocket",
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {},
	}
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-base.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.AddCommand(serveCmd)
	RootCmd.AddCommand(migrateCmd)
	RootCmd.AddCommand(workerCmd)
	RootCmd.AddCommand(kafkaCmd)
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
