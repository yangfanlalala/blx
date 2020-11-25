package wechat

import (
	"blxee.com/utils/crypto"
	"encoding/base64"
	"encoding/json"
)

func Decrypt(cipherText, iv, session string, obj interface{}) error {
	cipherDecode, _ := base64.StdEncoding.DecodeString(cipherText)
	ivDecode, _ := base64.StdEncoding.DecodeString(iv)
	sessionDecode, _ := base64.StdEncoding.DecodeString(session)
	plainText, err := crypto.AesDecrypt(cipherDecode, sessionDecode, ivDecode)
	if err != nil {
		return err
	}
	//fmt.Println(string(plainText))
	if err = json.Unmarshal(plainText, obj); err != nil {
		return err
	}
	return nil
}
