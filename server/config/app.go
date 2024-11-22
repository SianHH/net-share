package config

type App struct {
	Addr string `yaml:"addr"`
	Mode string `yaml:"mode"`

	Account  string `yaml:"account"`
	Password string `yaml:"password"`
	JwtKey   string `yaml:"jwt-key"`

	Ip          string   `yaml:"ip"`
	Domain      string   `json:"domain"`
	HostPort    string   `yaml:"host-port"`
	Entrypoint  string   `yaml:"entrypoint"`
	Ports       []string `yaml:"ports"`
	ForwardPort string   `yaml:"forward-port"`

	Logger Logger `yaml:"logger"`
}
