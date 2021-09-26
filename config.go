package main

import (
	"github.com/spf13/viper"
	"log"
)

var (
	Conf ConfigData
)

type ConfigData struct {
	MQTT_HOST string
	MQTT_PORT int

	MESSAGE_RATE     int
	MESSAGE_DURATION int

	DEVICE_STEP_NUM       int
	DEVICE_TRIAD_FILEPATH string
}

func loadConfig() {
	mustLoad(configEnv, configPath)
}

func mustLoad(env, filePath string) {
	v := viper.New()
	v.SetConfigFile(filePath)
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("ReadInConfig: %s", err)
	}
	v = v.Sub(env)
	err = v.UnmarshalExact(&Conf)
	if err != nil {
		log.Fatalf("UnmarshalExact: %s", err)
	}
}