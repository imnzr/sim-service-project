package helper

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func generateDLocalSignature(secretKey, method, path string, body []byte, xDate string) string {
	message := method + "\n" + path + "\n" + string(body) + "\n" + xDate
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(message))

	return hex.EncodeToString(mac.Sum(nil))
}
