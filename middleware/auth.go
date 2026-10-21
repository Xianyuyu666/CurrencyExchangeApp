package middleware

import (
	"ExchangeApp/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 检查Authorization是否有信息
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		// 检查格式：必须刚好两段且第一段是"Bearer"
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}
		// 取出token
		tokenStr := parts[1]

		// 验证token是否过期，篡改，格式错误
		jwtclaims, err := utils.ParseToken(tokenStr)

		if err != nil {
			log.Println("ParseToken error:", err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// 取出载荷信息,放行
		ctx.Set("UserID", jwtclaims.UserID)
		ctx.Next()
	}
}
