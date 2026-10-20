package controllers

import (
	"ExchangeApp/global"
	"ExchangeApp/models"
	"ExchangeApp/utils"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Req struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register 用户注册HandlerFunc
func Register(ctx *gin.Context) {
	var req Req
	// 获取请求体信息
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("ShouldBindJSON error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "server error"})
		return
	}

	var user models.User
	// 查询用户是否存在,同时用占位符防止SQL注入
	err := global.Db.Where("username = ?", req.Username).First(&user).Error

	// 用户不存在，可以创建
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Bcrypt加密
		hashedPwd, err := utils.HashPassword(req.Password)
		if err != nil {
			log.Println("HashPassword error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}

		// 数据库创建用户
		user = models.User{Username: req.Username, Password: hashedPwd}
		if err := global.Db.Create(&user).Error; err != nil {
			log.Println("Db.Create error:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Create Success"})
		return
	}
	// 数据库出错
	if err != nil {
		log.Println("Db.Where error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	// 用户名被占用
	ctx.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
}

// Login 用户登录HandlerFunc
func Login(ctx *gin.Context) {
	var req Req
	// 获取请求体信息
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("ShouldBindJSON error:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "server error"})
		return
	}

	var user models.User
	// 查询用户是否存在
	err := global.Db.Where("username = ?", req.Username).First(&user).Error

	// 数据库出错
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Db.Where error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	// 用户不存在或者密码错误
	if errors.Is(err, gorm.ErrRecordNotFound) || !utils.CheckPassword(user.Password, req.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// 用户存在且密码正确
	// 生成Token令牌给客户端，用于后续HTTP请求身份认证
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		log.Println("utils.GenerateJWT error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token, "message": "Login Success"})
}

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
