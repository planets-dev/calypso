# Calypso

Calypso is a small HTTPS reverse proxy written in Go. It terminates TLS with
Let's Encrypt certificates, supports HTTP/1.1, HTTP/2, and HTTP/3 per domain.

## Configuration

Create a TOML file with one `[[websites]]` entry per domain:

```toml
[[websites]]
domain = "example.com"
http1 = true
http2 = true
http3 = false
port = 5173
```
