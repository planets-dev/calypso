package main

import (
	"io"
	"net"
	"net/http"
	"strconv"
)

type application struct {
	httpClient    *http.Client
	hostToWebsite map[string]Website
}

var hopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}

func (app *application) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	host := req.Host
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	website, ok := app.hostToWebsite[host]
	if !ok {
		http.Error(w, "no such domain", http.StatusNotFound)
		return
	}

	url := req.URL.Clone()

	if website.Redirect != "" {
		url.Host = website.Redirect
		http.Redirect(w, req, url.String(), http.StatusMovedPermanently)
		return
	}

	url.Scheme = "http"
	url.Host = net.JoinHostPort("localhost", strconv.Itoa(int(website.Port)))

	proxyReq, err := http.NewRequestWithContext(req.Context(), req.Method, url.String(), req.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	proxyReq.Host = req.Host
	proxyReq.Header = req.Header.Clone()
	for _, header := range hopHeaders {
		proxyReq.Header.Del(header)
	}
	clientIP := req.RemoteAddr
	if host, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		clientIP = host
	}
	proxyReq.Header.Set("X-Real-IP", clientIP)

	proxyResp, err := app.httpClient.Do(proxyReq)
	if err != nil {
		http.Error(w, "upstream unavailable", http.StatusBadGateway)
		return
	}
	defer proxyResp.Body.Close()

	for key, values := range proxyResp.Header {
		hop := false
		for _, header := range hopHeaders {
			if http.CanonicalHeaderKey(key) == header {
				hop = true
				break
			}
		}
		if hop {
			continue
		}
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.Header().Set("Server", "calypso")

	if website.HTTP3 {
		w.Header().Set("Alt-Svc", `h3=":443"; ma=2592000`)
	}

	w.WriteHeader(proxyResp.StatusCode)
	_, _ = io.Copy(w, proxyResp.Body)
}
