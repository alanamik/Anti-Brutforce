package config

type Service struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Parameters struct {
	LimitLogin    int `yaml:"limitLogin"`
	LimitPassword int `yaml:"limitPassword"`
	LimitIP       int `yaml:"limitIp"`
}

type ListIPs struct {
	Path string `yaml:"pathFile"`
}
