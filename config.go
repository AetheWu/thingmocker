package thingmocker

import (
	"log"

	"github.com/spf13/viper"
)

var (
	Conf ConfigData = ConfigData{
		MQTT_HOST: "localhost",
		MQTT_PORT: 1883,

		MESSAGE_RATE:     1000,
		MESSAGE_DURATION: 3600 * 24 * 7,
		DEVICE_STEP_NUM:  100,
		DEVICE_NUM:       100,

		DEVICE_TRIAD_FILEPATH: "/etc/thingmocker/triads.csv",
		COMM_FILEPATH:         "/etc/thingmocker/comm.csv",
	}
)

type ConfigData struct {
	MQTT_HOST string `mapstructure:"mqtt_host"`
	MQTT_PORT int    `mapstructure:"mqtt_port"`
	MQTT_TLS  bool   `mapstructure:"mqtt_tls"`
	IF_ADDR   string `mapstructure:"if_addr"`

	MESSAGE_RATE     int `mapstructure:"message_rate"`
	MESSAGE_DURATION int `mapstructure:"message_duration"`

	DEVICE_STEP_NUM       int    `mapstructure:"device_step_num"`
	DEVICE_NUM            int    `mapstructure:"device_num"`
	DEVICE_TRIAD_FILEPATH string `mapstructure:"device_triad_filepath"`
	COMM_FILEPATH         string `mapstructure:"comm_filepath"`
}

func LoadConfig(env, path string) func() {
	return func() {
		mustLoad(env, path)
	}
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

	// v.BindPFlag("message_rate", mockCmd.Flag("rate"))
	// v.BindPFlag("message_duration", mockCmd.Flag("duration"))
	// v.BindPFlag("device_step_num", mockCmd.Flag("step"))
	// v.BindPFlag("device_num", mockCmd.Flag("num"))
	// v.BindPFlag("device_triad_filepath", mockCmd.Flag("device_triad_filepath"))
	// v.BindPFlag("if_addr", mockCmd.Flag("ifaddr"))

	err = v.Unmarshal(&Conf)
	if err != nil {
		log.Fatalf("UnmarshalExact: %s", err)
	}
}
