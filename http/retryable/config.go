package retryable

import (
	"net/url"
	"runtime"
	"time"
)

const (
	DialTimeout           = 5 * time.Second
	DialKeepAlive         = 30 * time.Second
	MaxIdleConnections    = 100
	IdleConnTimeout       = 90 * time.Second
	TLSHandshakeTimeout   = 10 * time.Second
	ExpectContinueTimeout = 10 * time.Second
)

type Options func(config *HttpClientConfig)

type HttpClientConfig struct {
	DialTimeout   time.Duration `json:"dial_timeout"`    // 建立TCP连接的超时时间
	DialKeepAlive time.Duration `json:"dial_keep_alive"` // TCP保活探测间隔

	MaxIdleConns          int           `json:"max_idle_conns"`          // 最大空闲连接数
	MaxIdleConnsPerHost   int           `json:"max_idle_conns_per_host"` // 每个主机的最大空闲连接数
	IdleConnTimeout       time.Duration `json:"idle_conn_timeout"`       // 空闲连接超时时间
	TLSHandshakeTimeout   time.Duration `json:"tls_handshake_timeout"`   // TLS握手超时时间
	ExpectContinueTimeout time.Duration `json:"expect_continue_timeout"` // 100-continue等待超时
	InsecureSkipVerify    bool          `json:"insecure_skip_verify"`    // 是否跳过TLS证书验证
	ProxyURL              *url.URL      `json:"proxy_url"`               // 代理服务器URL
}

func NewHttpClientConfig(opts ...Options) *HttpClientConfig {
	cliCfg := &HttpClientConfig{}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(cliCfg)
		}
	}

	return cliCfg
}

func WithDialTimeout(dialTimeout string) Options {
	return func(config *HttpClientConfig) {
		if d, err := time.ParseDuration(dialTimeout); err == nil {
			config.DialTimeout = d
		} else {
			config.DialTimeout = DialTimeout
		}
	}
}

func WithDialKeepAlive(dialKeepAlive string) Options {
	return func(config *HttpClientConfig) {
		if d, err := time.ParseDuration(dialKeepAlive); err == nil {
			config.DialKeepAlive = d
		} else {
			config.DialKeepAlive = DialKeepAlive
		}
	}
}

func WithMaxIdleConns(maxIdleConns int) Options {
	return func(config *HttpClientConfig) {
		if maxIdleConns > 0 {
			config.MaxIdleConns = maxIdleConns
		} else {
			config.MaxIdleConns = MaxIdleConnections
		}
	}
}

func WithMaxIdleConnsPerHost(maxIdleConnsPerHost int) Options {
	return func(config *HttpClientConfig) {
		if maxIdleConnsPerHost > 0 {
			config.MaxIdleConnsPerHost = maxIdleConnsPerHost
		} else {
			config.MaxIdleConnsPerHost = runtime.NumGoroutine()
		}
	}
}

func WithIdleConnTimeout(idleConnTimeout string) Options {
	return func(config *HttpClientConfig) {
		if d, err := time.ParseDuration(idleConnTimeout); err == nil {
			config.IdleConnTimeout = d
		} else {
			config.IdleConnTimeout = IdleConnTimeout
		}
	}
}

func WithTLSHandshakeTimeout(tlsHandshakeTimeout string) Options {
	return func(config *HttpClientConfig) {
		if d, err := time.ParseDuration(tlsHandshakeTimeout); err == nil {
			config.TLSHandshakeTimeout = d
		} else {
			config.TLSHandshakeTimeout = TLSHandshakeTimeout
		}
	}
}

func WithExpectContinueTimeout(expectContinueTimeout string) Options {
	return func(config *HttpClientConfig) {
		if d, err := time.ParseDuration(expectContinueTimeout); err == nil {
			config.ExpectContinueTimeout = d
		} else {
			config.ExpectContinueTimeout = ExpectContinueTimeout
		}
	}
}

func WithProxyURL(proxyURL string) Options {
	return func(config *HttpClientConfig) {
		if u, err := url.Parse(proxyURL); err == nil {
			config.ProxyURL = u
		}
	}
}
