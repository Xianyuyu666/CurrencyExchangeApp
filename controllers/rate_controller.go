package controllers

import (
	"ExchangeApp/global"
	"ExchangeApp/models"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type rateReq struct {
	From string  `json:"from_currency" binding:"required"`
	To   string  `json:"to_currency" binding:"required"`
	Rate float64 `json:"rate" binding:"required"`
	Date string  `json:"date" binding:"required"`
}

// CreateExchangeRate 插入汇率数据HandleFunc
func CreateExchangeRate(ctx *gin.Context) {
	var req rateReq
	// 获取请求体信息
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("ShouldBindJSON error:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// 货币类型统一大写，防止重复写入
	from := strings.ToUpper(req.From)
	to := strings.ToUpper(req.To)
	// 数据库查询三元组
	var rate models.ExchangeRate
	err := global.Db.Where(
		"from_currency = ? AND to_currency = ? AND date = ?",
		from, to, req.Date).First(&rate).Error

	// 汇率不存在，可以创建
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建模板
		rate = models.ExchangeRate{FromCurrency: from, ToCurrency: to, Date: req.Date, Rate: req.Rate}
		if err := global.Db.Create(&rate).Error; err != nil {
			log.Println("Db.Create error:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Create Success"})
		return
	}

	// 数据库出错
	if err != nil {
		log.Println("Db.Create error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	// 三元组已存在
	ctx.JSON(http.StatusConflict, gin.H{"error": "exchange rate already exists"})
}
