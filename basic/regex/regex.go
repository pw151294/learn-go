package main

import (
	"fmt"
	"regexp"
)

const (
	text    = "Redis server v=4.0.6 sha=00000000:0 malloc=jemalloc-4.0.3 bits=64 build=cddc4a876f0e5f1f" //nolint:gofumpt
	pattern = `\bv=(\d+)\.(\d+)\.(\d+)\b`
	cmdLine = "mysql --version"
)

func main() {
	re := regexp.MustCompile(pattern)
	fmt.Println(re.FindStringSubmatch(text))
}
