package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// rawUrl:http://172.29.230.157:5000/api/v0.1/ci_types/36/attributes
func buildQueryUrl(rawUrl string, params map[string]string) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", fmt.Errorf("failed to parse raw url %s: %v", rawUrl, err)
	}
	reqParams := url.Values{}
	for k, v := range params {
		reqParams.Add(k, v)
	}

	return fmt.Sprintf("%s?%s", u.String(), reqParams.Encode()), nil
}

func main() {
	req, err := http.NewRequest(http.MethodGet, "https://imooc.com/", nil)
	req.Header.Add("User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Println("redirect:", req)
			return nil
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	s, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(s))
}
