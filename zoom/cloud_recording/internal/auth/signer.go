package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/configs"
)

type Signer struct {
	cfg *configs.AuthConfig
}

func NewSigner(cfg *configs.AuthConfig) *Signer {
	return &Signer{cfg: cfg}
}

// GenerateToken 生成签名 token
// 算法：HMAC-SHA256(recordingID + "|" + expires, secret)，返回 Base64 URL 编码
func (s *Signer) GenerateToken(recordingID string, expires int64) string {
	message := fmt.Sprintf("%s|%d", recordingID, expires)
	mac := hmac.New(sha256.New, []byte(s.cfg.SignSecret))
	mac.Write([]byte(message))
	return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}

// GeneratePlayURL 生成完整播放 URL
// 格式：/play/{recordingID}?token={token}&expires={unix_timestamp}
func (s *Signer) GeneratePlayURL(recordingID string, expirySeconds int) string {
	expires := time.Now().Add(time.Duration(expirySeconds) * time.Second).Unix()
	token := s.GenerateToken(recordingID, expires)
	return fmt.Sprintf("/play/%s?token=%s&expires=%d", recordingID, token, expires)
}
