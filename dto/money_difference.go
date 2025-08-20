package dto

import (
	"math/big"
	"sort"
	"time"
)

// MoneyDifference representa la diferencia entre dos montos en fechas distintas
type MoneyDifference struct {
	FromDate         time.Time  `json:"fromDate"`
	ToDate           time.Time  `json:"toDate"`
	FromAmount       *big.Float `json:"fromAmount"`
	ToAmount         *big.Float `json:"toAmount"`
	Difference       *big.Float `json:"difference"`
	PercentageChange *big.Float `json:"percentageChange"`
}

// NewMoneyDifference crea una nueva instancia calculando diferencia y porcentaje
func NewMoneyDifference(fromDate, toDate time.Time, fromAmount, toAmount *big.Float) *MoneyDifference {
	diff := new(big.Float).Sub(toAmount, fromAmount)

	zero := big.NewFloat(0)
	var pctChange *big.Float
	if fromAmount.Cmp(zero) != 0 {
		pctChange = new(big.Float).Quo(diff, fromAmount)
		pctChange = pctChange.Mul(pctChange, big.NewFloat(100))
	} else {
		pctChange = big.NewFloat(0)
	}

	return &MoneyDifference{
		FromDate:         fromDate,
		ToDate:           toDate,
		FromAmount:       fromAmount,
		ToAmount:         toAmount,
		Difference:       diff,
		PercentageChange: pctChange,
	}
}

// Ordenar por fecha inicial ascendente
func SortByFromDateAsc(data []*MoneyDifference) {
	sort.Slice(data, func(i, j int) bool {
		return data[i].FromDate.Before(data[j].FromDate)
	})
}

// Ordenar por diferencia descendente
func SortByDifferenceDesc(data []*MoneyDifference) {
	sort.Slice(data, func(i, j int) bool {
		return data[i].Difference.Cmp(data[j].Difference) > 0
	})
}

// Ordenar por monto final ascendente
func SortByToAmountAsc(data []*MoneyDifference) {
	sort.Slice(data, func(i, j int) bool {
		return data[i].ToAmount.Cmp(data[j].ToAmount) < 0
	})
}

// Ordenar por cambio porcentual descendente
func SortByPercentageChangeDesc(data []*MoneyDifference) {
	sort.Slice(data, func(i, j int) bool {
		return data[i].PercentageChange.Cmp(data[j].PercentageChange) > 0
	})
}
