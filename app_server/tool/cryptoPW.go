package tool

import (
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/scrypt"
)

// ScryptyPw 密码加密
func ScryptyPw(password string) string {
	const KeyLen = 10
	salt := make([]byte, 8)
	// 盐
	salt = []byte{12, 32, 4, 6, 66, 33, 77, 88}

	HashPw, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen)
	if err != nil {
		logrus.Error("scrypt failed,err:", err)
		return ""
	}
	fpwd := base64.StdEncoding.EncodeToString(HashPw)
	return fpwd
}
