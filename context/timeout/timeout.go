package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc() //确保请求正常完成也会释放资源
	body := map[string]interface{}{
		"name":     "weipan4",
		"age":      24,
		"birthday": time.Now().Add(-4 * 365 * 24 * time.Hour),
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("marshal body fail: %s", err)
	}
	_, err = http.NewRequestWithContext(ctx, http.MethodGet,
		"http://localhost:8080", strings.NewReader(string(bytes))) //设置的ctx
	if err != nil {
		log.Fatalf("build new request fail: %s", err)
	}
}
