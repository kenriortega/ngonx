package proxy

// ProxyEndpoint struct for all enpoints
type ProxyEndpoint struct {
	Name      string     `mapstructure:"name"`
	HostURI   string     `mapstructure:"host_uri"`
	Endpoints []Endpoint `mapstructure:"endpoints"`
}

// Enpoint struct for enpoint object
type Endpoint struct {
	PathEndpoint  string `mapstructure:"path_endpoints"`
	PathToProxy   string `mapstructure:"path_proxy"`
	PathProtected bool   `mapstructure:"path_protected"`
}

// ProxyRepository interface
type ProxyRepository interface {
	SaveKEY(string, string, string) error
	GetKEY(string, string) (string, error)
}
