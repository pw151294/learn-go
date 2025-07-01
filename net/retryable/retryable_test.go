package retryable

import (
	"github.com/hashicorp/go-retryablehttp"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"net/http"
	"testing"
)

const filepath = "net/retryable/config.json"

func TestRetryableHttpClient(t *testing.T) {
	zap.InitLogger(zap.LogPath)
	InitRetryableHttpClient(filepath)

	req, err := http.NewRequest(http.MethodGet, "https://www.123.com/", nil)
	if err != nil {
		t.Fatal(err)
	}
	retryReq, err := retryablehttp.FromRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	RetryCli.Do(retryReq)
}
