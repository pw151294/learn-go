package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func openReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file faild: %w", err)
	}
	defer file.Close()
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file faild: %w", err)
	}
	return bytes, nil
}

func readConfig(fileName string) ([]byte, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get user home dir faild: %w", err)
	}

	configPath := filepath.Join(homeDir, fileName)
	buffer, err := openReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config faild: %w", err)
	}
	return buffer, nil
}

func main() {
	_, err := readConfig("setting.json")
	if err != nil {
		fmt.Println(err)
	}

	// 如果err是一个包装操作 使用errors.Unwrap 返回error中包含的下一个错误
	e := errors.Unwrap(err)
	fmt.Println(e)

	// 报告error链中是否包含和target值相同的错误
	if errors.Is(err, e) {
		fmt.Println(e)
	}
}
