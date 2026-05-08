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
	t.MaxIdleConnsPerHost = 500
	t.IdleConnTimeout = 60 * time.Second
	t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	t.ForceAttemptHTTP2 = false
	t.TLSNextProto = make(map[string]func(authority string, c *tls.Conn) http.RoundTripper)

	httpClient = &http.Client{
		Transport: t,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
