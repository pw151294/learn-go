package retryable

import (
	"encoding/json"
	"github.com/hashicorp/go-retryablehttp"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"log"
	"net/http"
	"os"
	"time"
)

const timeout = time.Second * 30

var RetryCli *retryablehttp.Client

func InitRetryableHttpClient(filepath string) {
	config, err := readConfig(filepath)
	if err != nil {
		log.Fatalf("read retryable net client config failed: %v", err)
	}

	cliCfg := NewHttpClientConfig(
		WithDialTimeout(config.DialTimeout),
		WithDialKeepAlive(config.DialKeepAlive),
		WithMaxIdleConns(config.MaxIdleConns),
		WithMaxIdleConnsPerHost(config.MaxIdleConnsPerHost),
		WithExpectContinueTimeout(config.ExpectContinueTimeout),
		WithIdleConnTimeout(config.IdleConnTimeout),
		WithTLSHandshakeTimeout(config.TlsHandshakeTimeout),
		WithMaxIdleConns(config.MaxIdleConns),
		WithProxyURL(config.ProxyUrl))

	RetryCli = NewRetryableClient(cliCfg, timeout)
}

func NewRetryableClient(clientConfig *HttpClientConfig, timeout time.Duration) *retryablehttp.Client {
	retryCli := retryablehttp.NewClient()

	retryCli.RetryWaitMin = time.Second
	retryCli.RetryWaitMax = 3 * time.Second
	retryCli.RetryMax = 3

	retryCli.RequestLogHook = func(logger retryablehttp.Logger, request *http.Request, retryNum int) { // 在执行重置之前执行
		zap.GetLogger().Info("begin retry ...", "retryNum", retryNum, "url", request.URL.String())
	}
	retryCli.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration { // 重试失败之后和下一次重试之间执行
		zap.GetLogger().Info("retry failed, waiting for next try", "retryNum", attemptNum)
		return min
	}

	retryCli.HTTPClient.Timeout = timeout
	retryCli.HTTPClient = newHttpClient(clientConfig)

	return retryCli
}

func readConfig(filepath string) (*Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	err = json.NewDecoder(file).Decode(config)
	return config, err
}
