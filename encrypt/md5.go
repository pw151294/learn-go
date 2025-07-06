package encrypt

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Encode(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
