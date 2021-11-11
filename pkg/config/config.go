package config

import (
	"fmt"
	"os"

	domain "github.com/kenriortega/ngonx/internal/proxy/domain"
	"github.com/kenriortega/ngonx/pkg/errors"
	"github.com/kenriortega/ngonx/pkg/logger"
	"github.com/spf13/viper"
)

// Config the main struct that define all elements
// inside `goproxy.yaml`
type Config struct {
	ProxyGateway `mapstructure:"proxy"`
	GrpcProxy    `mapstructure:"grpc"`
	StaticServer `mapstructure:"static_server"`
}

// GrpcProxy ...
type GrpcProxy struct {
	Listener       string         `mapstructure:"listener_grpc"`
	GrpcEndpoints  []GrpcEndpoint `mapstructure:"endpoints_grpc"`
	GrpcSSL        OptionSSL      `mapstructure:"ssl_grpc"`
	GrpcClientCert string         `mapstructure:"client_crt"`
}

type GrpcEndpoint struct {
	Name    string `mapstructure:"name"`
	HostURI string `mapstructure:"host_uri"`
}

// StaticServer struct for the static server obeject
type StaticServer struct {
	Host       string    `mapstructure:"host_server"`
	Port       int       `mapstructure:"port_server"`
	StaticFile string    `mapstructure:"static_files"`
	ServerSSL  OptionSSL `mapstructure:"ssl_server"`
}

// ProxyGateway struct for the proxy gateway object
type ProxyGateway struct {
	Host              string                 `mapstructure:"host_proxy"`
	Port              int                    `mapstructure:"port_proxy"`
	PortExporterProxy int                    `mapstructure:"port_exporter_proxy"`
	ProxySSL          OptionSSL              `mapstructure:"ssl_proxy"`
	ProxySecurity     ProxySecurity          `mapstructure:"security"`
	ProxyCache        ProxyCache             `mapstructure:"cache_proxy"`
	EnpointsProxy     []domain.ProxyEndpoint `mapstructure:"services_proxy"`
}

// OptionSSL struct for the ssl options
type OptionSSL struct {
	Enable  bool   `mapstructure:"enable"`
	SSLPort int    `mapstructure:"ssl_port"`
	CrtFile string `mapstructure:"crt_file"`
	KeyFile string `mapstructure:"key_file"`
}

// ProxySecurity struct for security object
type ProxySecurity struct {
	Type string `mapstructure:"type"`
}

// ProxyCache struct for cache options object
type ProxyCache struct {
	Engine string `mapstructure:"engine"`
	Key    string `mapstructure:"key"`
}

// LoadConfig load the config file from `path` and `name`
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

// CreateSettingFile create a setting file if it doesn`t exits
func (c *Config) CreateSettingFile(setingFile string) {
	f, err := os.Create(fmt.Sprintf("./%s", setingFile))
	ymldata :=
		`
static_server:
  host_server: 0.0.0.0
  port_server: 8080
  static_files: ./examples/dist
  ssl_server:
    enable: true
    ssl_port: 8443
    crt_file: ./key/cert.pem
    key_file: ./key/key.pem
grpc:
  listener_grpc: "0.0.0.0:50000"
  ssl_grpc:
  	enable: false
  	ssl_port: 50443
  	crt_file: ./ssl/cert.pem
  	key_file: ./ssl/key.pem
  endpoints_grpc:
  	- name: backend-1
  	- host_uri: 0.0.0.0:50050
proxy:
  host_proxy: 0.0.0.0
  port_proxy: 30000
  port_exporter_proxy: 10000
  ssl_proxy:
    enable: true
    ssl_port: 443
    crt_file: ./key/server.crt
    key_file: ./key/server.key
  cache_proxy:
    engine: badger # local|badgerDB|redis
    key: secretKey
  security:
    type: apikey # apikey|jwt|none
  # maps of microservices with routes
  services_proxy:
      - name: microA
        host_uri: http://localhost:3000
        endpoints:
          - path_endpoints: /api/v1/health/
            path_proxy: /health/
            path_protected: false

          - path_endpoints: /api/v1/version/
            path_proxy: /version/
            path_protected: true
`
	if err != nil {
		logger.LogError(errors.Errorf("config: %v", errors.ErrCreatingSettingFile).Error())

	}

	// defer f.Close()

	data := []byte(ymldata)

	_, err = f.Write(data)

	if err != nil {
		logger.LogError(errors.Errorf("config: %v", errors.ErrWritingSettingFile).Error())

	}
}
