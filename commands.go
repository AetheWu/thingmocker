package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	version = "0.0.1"
)

var (
	rootCmd = &cobra.Command{
		Use:   "thingmocker",
		Short: "Thing mockers for iot load-testing",
	}
	startCmd = &cobra.Command{
		Use:   "start [OPTIONS]",
		Short: "start mocking things",
		Run: func(cmd *cobra.Command, args []string) {
			StartMocker(Conf.DEVICE_TRIAD_FILEPATH, Conf.DEVICE_STEP_NUM, Conf.MESSAGE_RATE, Conf.MESSAGE_DURATION)
		},
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Thingmocker v%s\n", version)
		},
	}
)

var (
	configEnv  string
	configPath string

	deviceTriadFilePath string
	rate                int
	duration            int
	thingsAddStep       int
)

func init() {
	startCmd.Flags().StringVarP(&configEnv, "env", "e", "development", "config env")
	startCmd.Flags().StringVarP(&configPath, "config", "c", "/etc/thingmocker/config.yaml", "config file path")
	startCmd.Flags().IntVarP(&thingsAddStep, "step", "s", 1000, "num of new things added to iot platform per second")
	startCmd.Flags().StringVarP(&deviceTriadFilePath, "deviceTriadFilePath", "f", "/etc/thingmocker/triad.csv", "triad file to import")
	startCmd.Flags().IntVarP(&rate, "rate", "r", 1000, "message rate of things uploaded msg to iot platform")
	startCmd.Flags().IntVarP(&duration, "duration", "d", 60, "the time of load-testing")

	rootCmd.AddCommand(startCmd, versionCmd)
}

func Execute() error {
	cobra.OnInitialize(loadConfig)
	return rootCmd.Execute()
}
