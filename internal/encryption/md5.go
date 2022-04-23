package encryption

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
)

func GetMD5(token, password string) string {
	h := hmac.New(md5.New, []byte(token))
	h.Write([]byte(password))

	return hex.EncodeToString(h.Sum(nil))
}
