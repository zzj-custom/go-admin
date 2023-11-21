package initializer

import (
	"github.com/spf13/viper"
	"log"
)

func Config() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("Failed to read configuration file")
	}
}
