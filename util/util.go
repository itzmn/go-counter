package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

func Md5Sum(str string) string {

	m := md5.New()
	m.Write([]byte(str))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

func GenReqId() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5Sum(base64.URLEncoding.EncodeToString(b))
}
