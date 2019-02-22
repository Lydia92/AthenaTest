package encrypt

import (
	"encoding/hex"
	"github.com/astaxie/beego"
	"strings"
)

var encryptKey = beego.AppConfig.String("encrypt_key::key")

func EncryptStr(str string) string {
	cryptoStr, _ := AESEncryptWithECB([]byte(str), []byte(encryptKey))
	return strings.ToUpper(hex.EncodeToString(cryptoStr))
}

func DecryptStr(cryptSTr string) string {
	cryptStr2, _ := hex.DecodeString(cryptSTr)
	str, _ := AESDecryptWithECB(cryptStr2, []byte(encryptKey))
	return string(str)
}
