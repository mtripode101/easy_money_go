package model_test

import (
	"testing"

	"github.com/mtripode101/easymoney-go/model"
)

func TestParseMoney(t *testing.T) {
	em := model.EasyMoney{MoneyRaw: "$123.45"}
	amount, err := em.ParseMoney()
	if err != nil {
		t.Fatalf("ParseMoney failed: %v", err)
	}
	if amount != 123.45 {
		t.Errorf("Expected 123.45, got %f", amount)
	}
}

func TestParseMoneyDetailed(t *testing.T) {
	cases := []struct {
		raw      string
		expected float64
		currency string
	}{
		{"$123.45", 123.45, "$"},
		{"USD 123.45", 123.45, "USD"},
		{"€99,99", 99.99, "€"},
		{"ARS 1.234,56", 1234.56, "ARS"},
	}

	for _, c := range cases {
		em := model.EasyMoney{MoneyRaw: c.raw}
		amount, currency, err := em.ParseMoneyDetailed()
		if err != nil {
			t.Errorf("Error parsing %s: %v", c.raw, err)
		}
		if amount != c.expected || currency != c.currency {
			t.Errorf("Expected %f %s, got %f %s", c.expected, c.currency, amount, currency)
		}
	}
}
