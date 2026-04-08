package dingtalk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func genSign(secret string, timestamp int64) (string, error) {
	// stringToSign
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret

	// HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(stringToSign))
	if err != nil {
		return "", err
	}

	// base64
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signature, nil
}
