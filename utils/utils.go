package utils

import (
	"ExchangeApp/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

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

// GenerateJWT 生成JWT Token密钥：用于HTTP请求身份认证
func GenerateJWT(username string) (string, error) {
	// 填入Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		Username:  username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.AppConfig.Jwt.ExpireHours) * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	// 生成带签名Token
	signedToken, err := token.SignedString([]byte(config.AppConfig.Jwt.Secret))
	return "Bearer " + signedToken, err
}

// ParseToken 解析并验证JWT，成功返回*JWTClaims
func ParseToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 类型断言检查JWT算法是否正确
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		// 传入解析用的密钥
		return []byte(config.AppConfig.Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	// 类型断言检查token.Claims接口底层是否为JWTClaims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
