package controllers

import (
	"ExchangeApp/global"
	"ExchangeApp/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type articleReq struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required"`
}

// CreateArticle 发布文章（需登录）
func CreateArticle(ctx *gin.Context) {
	// 取出鉴权得到的用户ID
	userID, exists := ctx.Get("UserID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 类型断言检查
	uid, ok := userID.(uint)
	if !ok {
		log.Println("error: uid is not uint")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	// 绑定请求体
	var req articleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// 数据库插入数据
	article := models.Article{Title: req.Title, Content: req.Content, UserID: uid}
	if err := global.Db.Create(&article).Error; err != nil {
		log.Println("Db.Create error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	// 创建成功,返回创建结果
	ctx.JSON(http.StatusOK, gin.H{"message": "Create Success", "id": article.ID})
}
