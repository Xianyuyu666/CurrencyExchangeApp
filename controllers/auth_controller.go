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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// 查询用户是否存在,同时用占位符防止SQL注入
	err := global.Db.Where("username = ?", req.Username).First(&user).Error

	// 用户不存在，可以创建
	if errors.Is(err, gorm.ErrRecordNotFound){
		user = models.User{Username: req.Username, Password: req.Password}
		if err := global.Db.Create(&user).Error; err != nil {
			log.Println("Db.Create error:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Create Success"})
		return
	}
	// 数据库出错
	if err != nil {
		log.Println("Db.Where error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// 查询用户是否存在
	err := global.Db.Where("username = ?", req.Username).First(&user).Error

	// 数据库出错
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Db.Where error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 用户不存在或者密码错误
	if errors.Is(err, gorm.ErrRecordNotFound) || user.Password != req.Password {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}
	// 用户存在且密码正确
	ctx.JSON(http.StatusOK, gin.H{"message": "Login Success"})
}
