package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	APIUrl1 string `mapstructure:"API_URL_1"`
	APIUrl2 string `mapstructure:"API_URL_2"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	var config Config
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file: %s", err)
		return nil, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Error while unmarshalling config: %s", err)
		return nil, err
	}
	return &config, nil
}
