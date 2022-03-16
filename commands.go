package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			StartMocker(Conf.IF_ADDR, Conf.DEVICE_TRIAD_FILEPATH, Conf.DEVICE_STEP_NUM, Conf.MESSAGE_RATE, Conf.MESSAGE_DURATION, Conf.DEVICE_NUM)
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
	mockCmd.Flags().StringVarP(&configEnv, "env", "e", "development", "config env")
	mockCmd.Flags().StringVarP(&configPath, "config", "c", "/app/thingmocker/config.yaml", "config file path")

	mockCmd.Flags().IntVarP(&thingsAddStep, "step", "s", 1000, "num of new things added to iot platform per second")
	mockCmd.Flags().StringVarP(&deviceTriadFilePath, "device_triad_filepath", "f", "/etc/thingmocker/triad.csv", "triad file to import")
	mockCmd.Flags().IntVarP(&rate, "rate", "r", 1000, "message rate of things uploaded msg to iot platform")
	mockCmd.Flags().IntVarP(&duration, "duration", "d", 60, "the time of load-testing")
	mockCmd.Flags().IntVarP(&num, "num", "n", 100, "the number of things")
	mockCmd.Flags().StringVarP(&ifaddr, "ifaddr", "i", "", "interface addr")

	rootCmd.AddCommand(mockCmd, versionCmd)
}

func Execute() error {
	cobra.OnInitialize(loadConfig)
	return rootCmd.Execute()
}

var (
	Conf ConfigData
)

type ConfigData struct {
	MQTT_HOST string `mapstructure:"mqtt_host"`
	MQTT_PORT int    `mapstructure:"mqtt_port"`
	IF_ADDR   string `mapstructure:"if_addr"`

	MESSAGE_RATE     int `mapstructure:"message_rate"`
	MESSAGE_DURATION int `mapstructure:"message_duration"`

	DEVICE_STEP_NUM       int    `mapstructure:"device_step_num"`
	DEVICE_NUM            int    `mapstructure:"device_num"`
	DEVICE_TRIAD_FILEPATH string `mapstructure:"device_triad_filepath"`
}

func loadConfig() {
	mustLoad(configEnv, configPath)
}

func mustLoad(env, filePath string) {
	v := viper.New()
	v.SetConfigFile(filePath)

	v.SetDefault("author", "zhw")
	v.SetDefault("license", "apache")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("ReadInConfig: %s", err)
	}
	v = v.Sub(env)

	v.BindPFlag("message_rate", mockCmd.Flag("rate"))
	v.BindPFlag("message_duration", mockCmd.Flag("duration"))
	v.BindPFlag("device_step_num", mockCmd.Flag("step"))
	v.BindPFlag("device_num", mockCmd.Flag("num"))
	v.BindPFlag("device_triad_filepath", mockCmd.Flag("device_triad_filepath"))
	v.BindPFlag("if_addr", mockCmd.Flag("ifaddr"))

	err = v.UnmarshalExact(&Conf)
	if err != nil {
		log.Fatalf("UnmarshalExact: %s", err)
	}
}
