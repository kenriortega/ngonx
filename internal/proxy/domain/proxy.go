package proxy

// EndpointService ...
type ProxyEndpoint struct {
	Name      string     `mapstructure:"name"`
	HostURI   string     `mapstructure:"host_uri"`
	Endpoints []Endpoint `mapstructure:"endpoints"`
}

type Endpoint struct {
	PathEndpoint  string `mapstructure:"path_endpoints"`
	PathToProxy   string `mapstructure:"path_proxy"`
	PathProtected bool   `mapstructure:"path_protected"`
}

type ProxyRepository interface {
	SaveKEY(string, string, string) error
	GetKEY(string) (string, error)
}
