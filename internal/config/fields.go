package config

type Service struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Redis struct {
	Address string `yaml:"address"`
	DB      int    `yaml:"db"`
}

type Parameters struct {
	LimitLogin    int `yaml:"limitLogin"`
	LimitPassword int `yaml:"limitPassword"`
	LimitIP       int `yaml:"limitIP"`
}
