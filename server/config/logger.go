package config

type Logger struct {
	File    string `yaml:"file"`
	Level   int    `yaml:"level"`
	Console bool   `yaml:"console"`
}
