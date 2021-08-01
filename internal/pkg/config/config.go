package config

import (
	"fmt"
	"os"

	"github.com/kenriortega/goproxy/internal/pkg/errors"
	"github.com/kenriortega/goproxy/internal/pkg/logger"
	domain "github.com/kenriortega/goproxy/internal/proxy/domain"
	"github.com/spf13/viper"
)

type Config struct {
	ProxyGateway `mapstructure:"proxy"`
	StaticServer `mapstructure:"static_server"`
}

type StaticServer struct {
	Host       string    `mapstructure:"host_server"`
	Port       int       `mapstructure:"port_server"`
	StaticFile string    `mapstructure:"static_files"`
	ServerSSL  OptionSSL `mapstructure:"ssl_server"`
}

type ProxyGateway struct {
	Host          string                 `mapstructure:"host_proxy"`
	Port          int                    `mapstructure:"port_proxy"`
	ProxySSL      OptionSSL              `mapstructure:"ssl_proxy"`
	ProxySecurity ProxySecurity          `mapstructure:"security"`
	ProxyCache    ProxyCache             `mapstructure:"cache_proxy"`
	EnpointsProxy []domain.ProxyEndpoint `mapstructure:"services_proxy"`
}

type OptionSSL struct {
	Enable  bool   `mapstructure:"enable"`
	SSLPort int    `mapstructure:"ssl_port"`
	CrtFile string `mapstructure:"crt_file"`
	KeyFile string `mapstructure:"key_file"`
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
		err = errors.ErrReadConfig
		return
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		err = errors.ErrUnmarshalConfig
		return
	}
	return
}

func (c *Config) CreateSettingFile(setingFile string) {
	f, err := os.Create(fmt.Sprintf("./%s", setingFile))
	ymldata := `
proxy:
  host_proxy: 0.0.0.0
  port_proxy: 5000
  cache_proxy:
    engine: badger # local|badgerDB|redis
    key: secretKey
  security:
    type: jwt # apikey|jwt|none
    secret_key: key00 # apikey jwtkey this value can be replace by genkey command
  # maps of microservices with routes
  services_proxy:
      - name: microA
        host_uri: http://localhost:3000
        endpoints:
          - path_endpoints: /api/v1/health/
            path_proxy: /health/
            path_protected: false
`
	if err != nil {
		logger.LogError(errors.ErrCreatingSettingFile.Error())
	}

	defer f.Close()

	data := []byte(ymldata)

	_, err = f.Write(data)

	if err != nil {
		logger.LogError(errors.ErrWritingSettingFile.Error())
	}
}
