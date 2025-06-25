package retryable

import (
	"encoding/json"
	"flag"
	"github.com/hashicorp/go-retryablehttp"
	"log"
	"net/url"
	"os"
	"time"
)

const timeout = time.Second * 30

var filepath = flag.String("filepath", "http/retryable/config.json", "retryable client configuration file")

type Options struct {
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

func NewCli() *retryablehttp.Client {
	opts, err := readConfig(*filepath)
	if err != nil {
		log.Fatal(err)
	}
	return newRetryableClient(opts, timeout)
}

func newRetryableClient(opts *Options, timeout time.Duration) *retryablehttp.Client {
	retryCli := retryablehttp.NewClient()

	retryCli.HTTPClient.Timeout = timeout
	retryCli.HTTPClient = newHttpClient(opts)

	return retryCli
}

func readConfig(filepath string) (*Options, error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("failed to open retryable http client config file: %v", err)
	}
	defer file.Close()

	config := &Options{}
	if err = json.NewDecoder(file).Decode(config); err != nil {
		log.Fatal("load retryable http client config failed:", err)
	}

	return config, nil
}
