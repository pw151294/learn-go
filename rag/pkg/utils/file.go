package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// SaveFile 将数据写入指定路径的文件
func SaveFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// ReadFile 读取指定路径的文件内容
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// FileMD5 计算指定路径文件的MD5值
func FileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// SaveFileFromReader 将 *os.File 的内容保存到指定路径
func SaveFileFromReader(path string, data *os.File) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, data)
	return err
}
