package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const file = "/usr/local/zookeeper/bin/../conf/zoo.cfg"

func PrintFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	PrintFileContents(file)
}

func PrintFileContents(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func Exists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

// WalkDir 遍历指定目录下的所有文件（包括子目录）
// root: 要遍历的根目录路径
// 返回: 文件路径列表和错误信息
func WalkDir(root string) ([]string, error) {
	var files []string

	// 使用filepath.Walk遍历目录
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录，只收集文件
		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// Walk 收集指定目录下的所有文件
func Walk(root string) []string {
	fns := make([]string, 0)
	f := func(fp string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return err
		}
		fn := fi.Name()
		fns = append(fns, fn)
		return nil
	}
	if err := filepath.Walk(root, f); err != nil {
		log.Fatal(err)
	}

	return fns
}

func main() {
	//content := `1234'1234'"1234"`
	//PrintFileContents(strings.NewReader(content))
	//PrintFile(file)
	//log.Println(Exists(file))

	log.Println(filepath.Clean(file))
	log.Println(Walk("/usr/local/zookeeper"))
}
