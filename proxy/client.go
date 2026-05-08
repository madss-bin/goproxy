package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

var httpClient *http.Client

func init() {
	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     false,
		MaxIdleConns:          2000,
		MaxIdleConnsPerHost:   500,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second, // Prevents goroutine leaks if CDN hangs
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}

	httpClient = &http.Client{
		Transport: t,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
