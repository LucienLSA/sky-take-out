package utils

import (
	"crypto/md5"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordCost = 0 //密码加密难度
)

// 设置密码
func SetPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), PasswordCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// 检验密码
func CheckPassword(hashedpwd, pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedpwd), []byte(pwd))
}

// MD5V 对目标字符串取Hash salt：加盐字段，iteration：hash迭代轮数。
func MD5V(str string, salt string, iteration int) string {
	b := []byte(str)
	s := []byte(salt)
	h := md5.New()
	h.Write(s) // 先传入盐值，之前因为顺序错了卡了很久
	h.Write(b)
	var res []byte
	res = h.Sum(nil)
	for i := 0; i < iteration-1; i++ {
		h.Reset()
		h.Write(res)
		res = h.Sum(nil)
	}
	return hex.EncodeToString(res)
}
