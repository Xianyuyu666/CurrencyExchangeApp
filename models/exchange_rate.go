package models

import (
	"gorm.io/gorm"
)

type ExchangeRate struct {
	gorm.Model
	FromCurrency string  `gorm:"type:varchar(10);uniqueIndex:idx_from_to_date" json:"from_currency"`
	ToCurrency   string  `gorm:"type:varchar(10);uniqueIndex:idx_from_to_date" json:"to_currency"`
	Rate         float64 `json:"rate"`
	Date         string  `gorm:"type:varchar(10);uniqueIndex:idx_from_to_date" json:"date"`
}
