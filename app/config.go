package app

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Hashcat struct {
		Wordlist string `yaml:"wordlist"`
		Limit    int `yaml:"limit"`
	} `yaml:"hashcat"`
}

func ReadConfig() (*Config, error){
	f, err := os.Open("config.yaml")
	if err != nil {
		log.Println("Failed open config")
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Println("Failed read config")
		return nil, err
	}
	return &cfg, nil
}
