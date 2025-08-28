package encrypt

import (
	"encoding/json"
	"testing"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const encodedResp = "eyJjb2RlIjowLCJtc2ciOiJTVUNDRVNTIiwiZGF0YSI6eyIwXzE3Mi4zMC4zNC43MSI6MTczMX19"

func TestBase64Decode(t *testing.T) {
	respBody, err := Base64Decode([]byte(encodedResp))
	if err != nil {
		t.Fatal(err)
	}
	var result Result
	if err = json.Unmarshal([]byte(respBody), &result); err != nil {
		t.Fatal(err)
	}

	_, ok := result.Data.(map[string]int64)
	if !ok {
		t.Errorf("result data is not map[string]int64")
	}
}
