package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

func MD5(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

// Base64Encode 对字符串进行Base64编码
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64URLEncode URL安全的Base64编码
func Base64URLEncode(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

// Base64URLDecode URL安全的Base64解码
func Base64URLDecode(encodedStr string) (string, error) {
	decodedBytes, err := base64.URLEncoding.DecodeString(encodedStr)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

// Base64Decode 对Base64字符串进行解码
func Base64Decode(encodedStr string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedStr)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

func main() {

}
