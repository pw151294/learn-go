package retryable

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

const DefaultKeepAlive = 90 * time.Second

func newHttpClient(clientConfig *HttpClientConfig) *http.Client {
	if clientConfig == nil {
		return &http.Client{
			Transport: defaultTransport(),
		}
	}

	return &http.Client{
		Transport: newCliTransport(clientConfig),
	}
}

type Config struct {
	DialTimeout           string `json:"dial_timeout"`
	DialKeepAlive         string `json:"dial_keep_alive"`
	MaxIdleConns          int    `json:"max_idle_conns"`
	MaxIdleConnsPerHost   int    `json:"max_idle_conns_per_host"`
	IdleConnTimeout       string `json:"idle_conn_timeout"`
	TlsHandshakeTimeout   string `json:"tls_handshake_timeout"`
	ExpectContinueTimeout string `json:"expect_continue_timeout"`
	InsecureSkipVerify    bool   `json:"insecure_skip_verify"`
	ProxyUrl              string `json:"proxy_url"`
}

func defaultTransport() *http.Transport {
	return newCliTransport(defaultOptions())
}

func defaultOptions() *HttpClientConfig {
	return &HttpClientConfig{
		DialTimeout:           DialTimeout,
		DialKeepAlive:         DefaultKeepAlive,
		MaxIdleConns:          MaxIdleConnections,
		MaxIdleConnsPerHost:   runtime.NumGoroutine(),
		IdleConnTimeout:       IdleConnTimeout, // keep the same with keep-aliva
		TLSHandshakeTimeout:   TLSHandshakeTimeout,
		ExpectContinueTimeout: ExpectContinueTimeout,
	}
}

func newCliTransport(clientConfig *HttpClientConfig) *http.Transport {
	var (
		proxy func(*http.Request) (*url.URL, error)
	)

	if clientConfig.ProxyURL != nil {
		proxy = http.ProxyURL(clientConfig.ProxyURL)
	}

	return &http.Transport{
		Proxy: proxy,

		MaxIdleConns: func() int {
			if clientConfig.MaxIdleConns == 0 {
				return MaxIdleConnections
			}
			return clientConfig.MaxIdleConns
		}(),

		TLSClientConfig: func() *tls.Config {
			if clientConfig.InsecureSkipVerify {
				return &tls.Config{InsecureSkipVerify: true} //nolint:gosec
			}
			return &tls.Config{InsecureSkipVerify: false} //nolint:gosec
		}(),

		MaxIdleConnsPerHost: func() int {
			if clientConfig.MaxIdleConnsPerHost == 0 {
				return runtime.NumGoroutine()
			}
			return clientConfig.MaxIdleConnsPerHost
		}(),

		IdleConnTimeout: func() time.Duration {
			if clientConfig.IdleConnTimeout > time.Duration(0) {
				return clientConfig.IdleConnTimeout
			}
			return DefaultKeepAlive
		}(),

		TLSHandshakeTimeout: func() time.Duration {
			if clientConfig.TLSHandshakeTimeout > time.Duration(0) {
				return clientConfig.TLSHandshakeTimeout
			}
			return TLSHandshakeTimeout
		}(),

		ExpectContinueTimeout: func() time.Duration {
			if clientConfig.ExpectContinueTimeout > time.Duration(0) {
				return clientConfig.ExpectContinueTimeout
			}
			return ExpectContinueTimeout
		}(),
	}
}
