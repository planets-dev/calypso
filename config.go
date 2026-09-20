package main

import (
	"fmt"
	"io"

	"github.com/pelletier/go-toml/v2"
)

type Website struct {
	Domain   string `toml:"domain"`
	HTTP1    bool   `toml:"http1"`
	HTTP2    bool   `toml:"http2"`
	HTTP3    bool   `toml:"http3"`
	Port     uint16 `toml:"port"`
	Redirect string `toml:"redirect"`
}

type Config struct {
	Websites []Website `toml:"websites"`
}

func parseConfig(r io.Reader) (Config, error) {
	var cfg Config

	err := toml.NewDecoder(r).Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}
