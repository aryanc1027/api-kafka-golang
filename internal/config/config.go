package config 

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string `mapstructure:"server_address"`
	KafkaHosts []string `mapstructure:"kafka_hosts"`
	APIKey string `mapstructure:"api_key"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetDefault("server_address", ":8080")
	viper.SetDefault("kafka_hosts", []string{"localhost:9092"})
	viper.SetDefault("api_key", "your-default-api-key")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return &Config{
				ServerAddress: viper.GetString("server_address"),
				KafkaHosts: viper.GetStringSlice("kafka_hosts"),
				APIKey: viper.GetString("api_key"),
			}, nil
		}
		return nil, err
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil

}

