package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DataSource           string `yaml:"DataSource"`
	CacheSource          string `yaml:"CacheSource"`
	CaptchaToken         string `yaml:"CaptchaToken"`
	AccessLog            string `yaml:"AccessLog"`
	ErrorLog             string `yaml:"ErrorLog"`
	MaxConcurrency       int    `yaml:"MaxConcurrency"`
	ErrorSleepTime       int    `yaml:"ErrorSleepTime"`
	EverySolveSleepTime  int    `yaml:"EverySolveSleepTime"`
	DisPatcherTickerTime int    `yaml:"DisPatcherTickerTime"`
}

func NewConfig(cfg string) *Config {
	file, err := os.Open(cfg)
	if err != nil {
		panic(err)
	}
	bts, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var c Config
	if err := yaml.Unmarshal(bts, &c); err != nil {
		panic(err)
	}

	return &c
}
