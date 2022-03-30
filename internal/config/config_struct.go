package config

import "time"

type Configuration struct {
	Listen  Listen        `yaml:"listen,omitempty"`
	Quotes  Quotes        `yaml:"quotes,omitempty"`
	Cors    Cors          `yaml:"cors,omitempty"`
	Runtime RuntimeConfig `yaml:"runtime,omitempty"`
}

type Quotes struct {
	Url  string `yaml:"url,omitempty"`
	Path string `yaml:"path,omitempty"`
}

type RuntimeConfig struct {
	GOGC       int `yaml:"gogc,omitempty"`
	GOMAXPROCS int `yaml:"gomaxprocs"`
}

type Listen struct {
	Host string `yaml:"host,omitempty"`
	Port int    `yaml:"port,omitempty"`
}

type Cors struct {
	Enabled     bool          `yaml:"enabled"`
	Origins     []string      `yaml:"origins"`
	MaxAge      time.Duration `yaml:"maxage"`
	Credentials bool          `yaml:"credentials"`
	Allow       AllowConfig   `yaml:"allow"`
	Expose      ExposeConfig  `yaml:"expose"`
}

type AllowConfig struct {
	Methods []string `yaml:"methods"`
	Headers []string `yaml:"headers"`
}

type ExposeConfig struct {
	Headers []string `yaml:"headers"`
}
