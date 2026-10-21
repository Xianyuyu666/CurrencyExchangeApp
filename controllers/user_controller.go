package controllers

import (
	"ExchangeApp/global"
	"ExchangeApp/models"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Info 用户获取自身信息HandleFunc
func Info(ctx *gin.Context) {
	// 获取用户信息
	usernameVal, exists := ctx.Get("Username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 类型断言
	usernameStr, ok := usernameVal.(string)
	if !ok {
		log.Println("username type assertion failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	// 数据库查询取数据
	var user models.User
	err := global.Db.Where("username = ?", usernameStr).First(&user).Error

	// 数据库出错
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Db.Where error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	// 用户不存在
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username"})
		return
	}

	// 发送到客户端,只能发送安全字段，注意敏感字段
	ctx.JSON(http.StatusOK, gin.H{
		"username":  user.Username,
		"createdAt": user.CreatedAt,
	})
}
