package controllers

import (
	"ExchangeApp/global"
	"ExchangeApp/models"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

type articleListQuery struct {
	Keyword  string `form:"keyword"`
	Page     int    `form:"page"`
	PageSize int    `form:"pagesize"`
}

// ListArticles 分页展示文章列表，支持keyword模糊查询
func ListArticles(ctx *gin.Context) {
	// 获取URL查询参数
	var req articleListQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	// 数据库查询
	// 动态调节查询条件先取出表
	db := global.Db.Model(&models.Article{})
	if req.Keyword != "" {
		db = db.Where("title LIKE ?", "%"+req.Keyword+"%")
	}
	// 总数
	var total int64
	db.Count(&total)

	// 数据库查询
	var list []models.Article
	err := db.Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username")
	}).
		Limit(req.PageSize).
		Offset((req.Page - 1) * req.PageSize).
		Order("created_at DESC").
		Find(&list).Error

	if err != nil {
		log.Println("Db.Find error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	// 成功返回文章全部信息，关联的用户信息只有id和username，再返回一个查到文章总数
	ctx.JSON(http.StatusOK, gin.H{"list": list, "total": total})
}

// GetArticle 获取指定id的文章
func GetArticle(ctx *gin.Context) {
	// 从URL中获取信息
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// 查询
	var article models.Article
	err = global.Db.Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username")
	}).First(&article, id).Error

	// 查不到
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}
	// 数据库错误
	if err != nil {
		log.Println("Db.First error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	// 成功返回文章信息加部分作者信息
	ctx.JSON(http.StatusOK, gin.H{
		"id":         article.ID,
		"title":      article.Title,
		"content":    article.Content,
		"likecount":  article.LikeCount,
		"author":     article.Author,
		"created_at": article.CreatedAt,
	})
}
