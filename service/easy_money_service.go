package service

import (
	"errors"
	"math/big"
	"sort"
	"time"

	"github.com/mtripode101/easymoney-go/config"
	"github.com/mtripode101/easymoney-go/dto"
	"github.com/mtripode101/easymoney-go/model"
	"github.com/mtripode101/easymoney-go/repository"
)

// Helper para convertir MoneyRaw a *big.Float
func parseMoney(raw string) *big.Float {
	val, ok := new(big.Float).SetString(raw)
	if ok {
		return val
	}
	return big.NewFloat(0)
}

// Guardar nuevo registro
func Save(em *model.EasyMoney) error {
	return config.DB.Create(em).Error
}

// Actualizar registro existente
func Update(id uint, updated *model.EasyMoney) (*model.EasyMoney, error) {
	var existing model.EasyMoney
	if err := config.DB.First(&existing, id).Error; err != nil {
		return nil, err
	}

	existing.Date = updated.Date
	existing.Description = updated.Description
	existing.MoneyRaw = updated.MoneyRaw

	err := config.DB.Save(&existing).Error
	return &existing, err
}

// Eliminar por ID
func Delete(id uint) error {
	return config.DB.Delete(&model.EasyMoney{}, id).Error
}

// Buscar todos ordenados por fecha descendente
func FindAll() ([]model.EasyMoney, error) {
	var results []model.EasyMoney
	err := config.DB.Order("date DESC").Find(&results).Error
	return results, err
}

// Buscar por rango de fechas
func FindByDateRange(start, end time.Time) ([]model.EasyMoney, error) {
	return repository.FindByDateBetween(start, end)
}

// Buscar por descripción
func SearchByDescription(keyword string) ([]model.EasyMoney, error) {
	return repository.FindByDescriptionContaining(keyword)
}

// Buscar por ID
func FindByID(id uint) (*model.EasyMoney, error) {
	var em model.EasyMoney
	if err := config.DB.First(&em, id).Error; err != nil {
		return nil, errors.New("registro no encontrado con id")
	}
	return &em, nil
}

// Calcular suma total entre fechas
func CalculateDifference(start, end time.Time) (*big.Float, error) {
	entries, err := repository.FindByDateBetween(start, end)
	if err != nil {
		return nil, err
	}

	total := big.NewFloat(0)
	for _, e := range entries {
		total = total.Add(total, parseMoney(e.MoneyRaw))
	}
	return total, nil
}

// Calcular diferencias consecutivas
func CalculateConsecutiveDifferences(start, end time.Time) ([]*dto.MoneyDifference, error) {
	entries, err := repository.FindByDateBetween(start, end)
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Date.Before(entries[j].Date)
	})

	var differences []*dto.MoneyDifference
	for i := 1; i < len(entries); i++ {
		prev := entries[i-1]
		curr := entries[i]
		diff := dto.NewMoneyDifference(
			prev.Date,
			curr.Date,
			parseMoney(prev.MoneyRaw),
			parseMoney(curr.MoneyRaw),
		)
		differences = append(differences, diff)
	}

	return differences, nil
}

// Calcular diferencia total entre primera y última entrada
func CalculateTotalMoneyDifference(start, end time.Time) (*dto.MoneyDifference, error) {
	entries, err := repository.FindByDateBetween(start, end)
	if err != nil {
		return nil, err
	}

	if len(entries) < 2 {
		return nil, errors.New("se necesitan al menos dos registros para calcular la diferencia total")
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Date.Before(entries[j].Date)
	})

	first := entries[0]
	last := entries[len(entries)-1]

	return dto.NewMoneyDifference(
		first.Date,
		last.Date,
		parseMoney(first.MoneyRaw),
		parseMoney(last.MoneyRaw),
	), nil
}
