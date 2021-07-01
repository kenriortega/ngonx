package internal

import (
	domain "egosystem.org/micros/gateway/domain"
	"github.com/spf13/viper"
)

type Config struct {
	ProxyGateway `mapstructure:"proxy"`
}

type ProxyGateway struct {
	Host          string          `mapstructure:"host_proxy"`
	Port          int             `mapstructure:"port_proxy"`
	EnpointsProxy []EnpointsProxy `mapstructure:"enpoints_proxy"`
}
type EnpointsProxy struct {
	Name     string                   `mapstructure:"name"`
	Services []domain.EndpointService `mapstructure:"services"`
}

func LoadConfig(path, name string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
