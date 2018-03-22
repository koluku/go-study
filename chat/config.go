package main

import (
	"github.com/BurntSushi/toml"
)

// Config is
type Config struct {
	Discord OAuth2
}

// OAuth2 is
type OAuth2 struct {
	ID     string
	Secret string
}

func getConfig() *Config {
	var config Config
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		panic(err)
	}
	return &config
}
