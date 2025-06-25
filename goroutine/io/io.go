package main

import (
	"fmt"
	"os"
	"time"
)

const (
	file1 = "config/cmd/config/config.toml"
	file2 = "config/yaml/config.yaml"
	file3 = "config/toml/config.toml"
)

// 示例：并发读取多个文件
func readFile(filename string, results chan string) {
	file, err := os.ReadFile(filename)
	if err != nil {
		results <- fmt.Sprintf("read file failed, filename: %s, error: %v", filename, err)
	}
	results <- fmt.Sprintf("read file %s, %d bytes", filename, len(file))
}

func main() {
	files := []string{file1, file2, file3}
	resCh := make(chan string, len(files))
	for _, f := range files {
		go readFile(f, resCh)
	}

	time.Sleep(1 * time.Second)
	close(resCh)

	for res := range resCh { // 如果resCh没有关闭的话 这里会出现deadlock
		fmt.Println(res)
	}
}
