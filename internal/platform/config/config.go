package config

import (
	domain "github.com/kenriortega/goproxy/internal/proxy/domain"
	"github.com/spf13/viper"
)

type Config struct {
	ProxyGateway `mapstructure:"proxy"`
}

type ProxyGateway struct {
	Host          string                 `mapstructure:"host_proxy"`
	Port          int                    `mapstructure:"port_proxy"`
	ProxySecurity ProxySecurity          `mapstructure:"security"`
	ProxyCache    ProxyCache             `mapstructure:"cache_proxy"`
	EnpointsProxy []domain.ProxyEndpoint `mapstructure:"services_proxy"`
}

type ProxySecurity struct {
	Type string `mapstructure:"type"`
}
type ProxyCache struct {
	Engine string `mapstructure:"engine"`
	Key    string `mapstructure:"key"`
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
