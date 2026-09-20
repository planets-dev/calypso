package main

type Website struct {
	Domain string `toml:"domain"`
	HTTP1  bool   `toml:"http1"`
	HTTP2  bool   `toml:"http2"`
	HTTP3  bool   `toml:"http3"`
	Port   uint16 `toml:"port"`
}

type Config struct {
	Websites []Website `toml:"websites"`
}
