package httpclient

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

func NewHTTPClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = 10
	transport.GetProxyConnectHeader = func(ctx context.Context, proxyURL *url.URL, target string) (http.Header, error) {
		return http.Header{"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"}}, nil
	}
	cl := http.Client{
		Timeout:   time.Duration(10) * time.Second,
		Transport: transport,
	}
	return &cl
}
