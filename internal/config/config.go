package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Port int
	Url  string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) LoadConfig(vendor string) error {
	v := viper.New()
	v.SetDefault("apiServer.port", 8080)

	// read environment variables.  env vars have precedence over config file
	v.AutomaticEnv()

	// read config file
	v.AddConfigPath(".")
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	c.Port = v.GetInt("apiServer.port")
	c.Url = v.GetString("nwsClient.url")

	if !v.IsSet("apiServer.port") || c.Port == 0 {
		return fmt.Errorf("port not set")
	}
	if !v.IsSet("nwsClient.url") || c.Url == "" {
		return fmt.Errorf("URL not set")
	}

	return nil
}
