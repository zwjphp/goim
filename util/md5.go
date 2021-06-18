package util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func MakePasswd(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}
