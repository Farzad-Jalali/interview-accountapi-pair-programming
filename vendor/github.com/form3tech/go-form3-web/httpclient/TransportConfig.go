package httpclient

func DefaultTransportConfig() *TransportConfig {
	return &TransportConfig{
		BasePath: "/v1",
	}
}

type TransportConfig struct {
	Host     string
	BasePath string
	Schemes  []string
}

func (cfg *TransportConfig) WithHost(host string) *TransportConfig {
	cfg.Host = host
	return cfg
}

func (cfg *TransportConfig) WithBasePath(basePath string) *TransportConfig {
	cfg.BasePath = basePath
	return cfg
}

func (cfg *TransportConfig) WithSchemes(schemes []string) *TransportConfig {
	cfg.Schemes = schemes
	return cfg
}
