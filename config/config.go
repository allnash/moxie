package config

type Config struct {
	Version    string    `yaml:"version"`
	Host	   string    `yaml:"host"`
	Port       string    `yaml:"port"`
	Logfile     string    `yaml:"log_file"`
	Services   []Service `yaml:"services"`
}

type Service struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"` // Type is one of ['web', 'proxy', 'static']
	IngressUrl string `yaml:"ingress_url"`
	EgressUrl  string `yaml:"egress_url"`
}
