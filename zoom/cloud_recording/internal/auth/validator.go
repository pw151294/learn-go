package auth

import (
	"crypto/hmac"
	"fmt"
	"strconv"
	"time"

	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/configs"
)

type Validator struct {
	signer *Signer
}

func NewValidator(cfg *configs.AuthConfig) *Validator {
	return &Validator{signer: NewSigner(cfg)}
}

// ValidateToken 验证签名
// 1. 检查是否过期（expires < now）
// 2. 重新计算 token
// 3. 使用 hmac.Equal 对比（防止时序攻击）
func (v *Validator) ValidateToken(recordingID, token, expiresStr string) error {
	expires, err := strconv.ParseInt(expiresStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid expires: %w", err)
	}
	if time.Now().Unix() > expires {
		return fmt.Errorf("token expired")
	}
	expected := v.signer.GenerateToken(recordingID, expires)
	// 注意：需要将 token 和 expected 都转为 []byte 再用 hmac.Equal 比较
	if !hmac.Equal([]byte(token), []byte(expected)) {
		return fmt.Errorf("invalid token")
	}
	return nil
}
