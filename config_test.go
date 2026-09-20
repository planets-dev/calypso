package main

import (
	"strings"
	"testing"

	"github.com/planets-dev/calypso/internal/assert"
)

func TestParseConfig(t *testing.T) {
	cfg, err := parseConfig(strings.NewReader(`
[[websites]]
domain = "example.com"
http1 = true
http2 = true
http3 = false
port = 443
`))
	assert.Nil(t, err)
	assert.Equal(t, len(cfg.Websites), 1)
	assert.Equal(t, cfg.Websites[0], Website{
		Domain: "example.com",
		HTTP1:  true,
		HTTP2:  true,
		HTTP3:  false,
		Port:   443,
	})
}

func TestParseConfigRedirect(t *testing.T) {
	cfg, err := parseConfig(strings.NewReader(`
[[websites]]
domain = "example.com"
http1 = true
http2 = true
http3 = false
redirect = "https://example.com"
`))
	assert.Nil(t, err)
	assert.Equal(t, cfg.Websites[0], Website{
		Domain:   "example.com",
		HTTP1:    true,
		HTTP2:    true,
		HTTP3:    false,
		Redirect: "https://example.com",
	})
}

func TestParseConfigMultipleWebsites(t *testing.T) {
	cfg, err := parseConfig(strings.NewReader(`
[[websites]]
domain = "example.com"
http1 = true
http2 = true
http3 = false
port = 443

[[websites]]
domain = "redirect.example.com"
http1 = true
http2 = false
http3 = true
redirect = "https://example.com"
`))
	assert.Nil(t, err)
	assert.Equal(t, len(cfg.Websites), 2)
	assert.Equal(t, cfg.Websites[0], Website{Domain: "example.com", HTTP1: true, HTTP2: true, Port: 443})
	assert.Equal(t, cfg.Websites[1], Website{Domain: "redirect.example.com", HTTP1: true, HTTP3: true, Redirect: "https://example.com"})
}
