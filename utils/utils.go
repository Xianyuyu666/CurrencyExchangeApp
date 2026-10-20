package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 密码加密
func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	return string(hash), err
}

// CheckPassword 校验密码
// 自动提取盐，重新哈希比对
func CheckPassword(hashedPwd, plainPwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)) == nil
}
