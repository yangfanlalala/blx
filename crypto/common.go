package crypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func HmacSha256Hex(key, str string) string {
	hax := hmac.New(sha256.New, []byte(key))
	hax.Write([]byte(str))
	return hex.EncodeToString(hax.Sum(nil))
}

func Sha256(plain string) string {
	sha := sha256.New()
	sha.Write([]byte(plain))
	bs := sha.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func HmacSha1Hex(key, str string) string {
	hx := hmac.New(sha1.New, []byte(key))
	hx.Write([]byte(str))
	return fmt.Sprintf("%s", hx.Sum(nil))
}

func HmacSha1(str, key string) string {
	hx := hmac.New(sha1.New, []byte(key))
	hx.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(hx.Sum(nil))
}
