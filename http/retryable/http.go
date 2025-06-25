package retryable

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

const DefaultKeepAlive = 90 * time.Second

func newHttpClient(opt *Options) *http.Client {
	if opt == nil {
		return &http.Client{
			Transport: defaultTransport(),
		}
	}

	return &http.Client{
		Transport: newCliTransport(opt),
	}
}

func defaultTransport() *http.Transport {
	return newCliTransport(defaultOptions())
}

func defaultOptions() *Options {
	return &Options{
		DialTimeout:           30 * time.Second,
		DialKeepAlive:         DefaultKeepAlive,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   runtime.NumGoroutine(),
		IdleConnTimeout:       DefaultKeepAlive, // keep the same with keep-aliva
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
	}
}

func newCliTransport(opt *Options) *http.Transport {
	var (
		proxy func(*http.Request) (*url.URL, error)
	)

	if opt.ProxyURL != nil {
		proxy = http.ProxyURL(opt.ProxyURL)
	}

	return &http.Transport{
		Proxy: proxy,

		MaxIdleConns: func() int {
			if opt.MaxIdleConns == 0 {
				return 100
			}
			return opt.MaxIdleConns
		}(),

		TLSClientConfig: func() *tls.Config {
			if opt.InsecureSkipVerify {
				return &tls.Config{InsecureSkipVerify: true} //nolint:gosec
			}
			return &tls.Config{InsecureSkipVerify: false} //nolint:gosec
		}(),

		MaxIdleConnsPerHost: func() int {
			if opt.MaxIdleConnsPerHost == 0 {
				return runtime.NumGoroutine()
			}
			return opt.MaxIdleConnsPerHost
		}(),

		IdleConnTimeout: func() time.Duration {
			if opt.IdleConnTimeout > time.Duration(0) {
				return opt.IdleConnTimeout
			}
			return DefaultKeepAlive
		}(),

		TLSHandshakeTimeout: func() time.Duration {
			if opt.TLSHandshakeTimeout > time.Duration(0) {
				return opt.TLSHandshakeTimeout
			}
			return 10 * time.Second
		}(),

		ExpectContinueTimeout: func() time.Duration {
			if opt.ExpectContinueTimeout > time.Duration(0) {
				return opt.ExpectContinueTimeout
			}
			return time.Second
		}(),
	}
}
