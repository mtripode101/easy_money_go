package dto_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/mtripode101/easymoney-go/dto"
)

func newAmount(val float64) *big.Float {
	return big.NewFloat(val)
}

func TestNewMoneyDifference(t *testing.T) {
	from := newAmount(100)
	to := newAmount(150)
	md := dto.NewMoneyDifference(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC), from, to)

	if md.Difference.Cmp(newAmount(50)) != 0 {
		t.Errorf("Expected difference 50, got %v", md.Difference)
	}
	if md.PercentageChange.Cmp(newAmount(50)) != 0 {
		t.Errorf("Expected percentage change 50%%, got %v", md.PercentageChange)
	}
}
