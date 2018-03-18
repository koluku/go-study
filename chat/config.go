package main

import (
	"github.com/BurntSushi/toml"
)

// Config is
type Config struct {
	Oauth2 oauth2Config
}

type oauth2Config struct {
	ID     int
	Secret string
}

func getConfig() *Config {
	var config Config
	_, err := toml.DecodeFile("config.tml", &config)
	if err != nil {
		panic(err)
	}
	return &config
}
