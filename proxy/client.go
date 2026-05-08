package proxy

import (
	"crypto/tls"
	"net/http"
	"time"
)

var httpClient *http.Client

func init() {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 2000
	t.MaxConnsPerHost = 500
	t.MaxIdleConnsPerHost = 500
	t.IdleConnTimeout = 60 * time.Second
	t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	httpClient = &http.Client{
		Transport: t,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
