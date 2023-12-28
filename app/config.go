package app

import (
	"gopkg.in/yaml.v3"
	"os"
)

type AppConfig struct {
	ScrapeSeconds      int  `yaml:"scrape_seconds"`
	TimeoutSeconds     int  `yaml:"timeout_seconds"`
	IgnoreCertificates bool `yaml:"ignore_certificates,omitempty"`
	PrometheusPort     int  `yaml:"prom_port"`
}

func LoadAppConfig(appConfigFile string) (*AppConfig, error) {
	f, err := os.ReadFile(appConfigFile)
	if err != nil {
		return nil, err
	}

	var appConfig AppConfig
	err = yaml.Unmarshal(f, &appConfig)
	if err != nil {
		return nil, err
	}

	return &appConfig, nil
}
