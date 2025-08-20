package repository

import (
	"time"

	"github.com/mtripode101/easymoney-go/config"
	"github.com/mtripode101/easymoney-go/model"
)

// Buscar por rango de fechas
func FindByDateBetween(start, end time.Time) ([]model.EasyMoney, error) {
	var results []model.EasyMoney
	err := config.DB.Where("date BETWEEN ? AND ?", start, end).Find(&results).Error
	return results, err
}

// Buscar por palabra clave en descripción (case-insensitive)
func FindByDescriptionContaining(keyword string) ([]model.EasyMoney, error) {
	var results []model.EasyMoney
	pattern := "%" + keyword + "%"
	err := config.DB.Where("LOWER(description) LIKE LOWER(?)", pattern).Find(&results).Error
	return results, err
}
