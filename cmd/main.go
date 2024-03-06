package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"thingmocker"
)

const (
	version = "1.0.0"
)

var (
	rootCmd = &cobra.Command{
		Use:   "thingmocker",
		Short: "Thing mockers for iot load-testing",
	}
	mockCmd = &cobra.Command{
		Use:   "mock [OPTIONS]",
		Short: "mock mocking things",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("env: %s, load config: %s", configEnv, configPath)
			thingmocker.LoadConfig(configEnv, configPath)()
			err := thingmocker.Run()
			if err != nil {
				cmd.PrintErr(err)
			}
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
	num                 int
	ifaddr              string
)

func init() {
	mockCmd.PersistentFlags().StringVarP(&configEnv, "env", "e", "defaults", "config env")
	mockCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "/app/thingmocker/config.yaml", "config file path")

	mockCmd.Flags().IntVarP(&thingsAddStep, "step", "s", 1000, "num of new things added to iot platform per second")
	mockCmd.Flags().StringVarP(&deviceTriadFilePath, "device_triad_filepath", "f", "/etc/thingmocker/triad.csv", "triad file to import")
	mockCmd.Flags().IntVarP(&rate, "rate", "r", 1000, "message rate of things uploaded msg to iot platform")
	mockCmd.Flags().IntVarP(&duration, "duration", "d", 60, "the time of load-testing")
	mockCmd.Flags().IntVarP(&num, "num", "n", 100, "the number of things")
	mockCmd.Flags().StringVarP(&ifaddr, "ifaddr", "i", "", "interface addr")

	rootCmd.AddCommand(mockCmd, versionCmd)
}

func main() {
	Execute()
}

func Execute() error {
	return rootCmd.Execute()
}
