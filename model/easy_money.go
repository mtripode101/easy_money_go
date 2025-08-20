package model

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// EasyMoney representa una entrada de dinero con fecha y descripción.
// El campo MoneyRaw se guarda como texto en la base de datos (VARCHAR),
// y se usa directamente en formularios y vistas.
type EasyMoney struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Date        time.Time `json:"date" form:"Date" time_format:"2006-01-02"`
	MoneyRaw    string    `gorm:"column:money" json:"money" form:"MoneyRaw"`
	Description string    `json:"description" form:"Description"`
}

// Método para convertir MoneyRaw a float64
func (em *EasyMoney) ParseMoney() (float64, error) {
	clean := strings.ReplaceAll(em.MoneyRaw, "$", "")
	return strconv.ParseFloat(clean, 64)
}

// ParseMoneyDetailed interpreta MoneyRaw y devuelve monto, moneda y error.
// Soporta formatos como "$123.45", "USD 123.45", "€99,99", "ARS 1.234,56".
func (em *EasyMoney) ParseMoneyDetailed() (float64, string, error) {
	raw := strings.TrimSpace(em.MoneyRaw)
	if raw == "" {
		return 0, "", errors.New("empty money string")
	}

	var currency string
	var amountStr string

	// Detectar si el formato es tipo "USD 123.45"
	parts := strings.Fields(raw)
	if len(parts) == 2 {
		currency = parts[0]
		amountStr = parts[1]
	} else {
		// Ej: "$123.45" o "€99,99"
		runes := []rune(raw)
		if !unicode.IsDigit(runes[0]) {
			currency = string(runes[0])
			amountStr = string(runes[1:])
		} else {
			return 0, "", errors.New("invalid format: missing currency symbol or code")
		}
	}

	amountStr = strings.TrimSpace(amountStr)

	// Normalizar separadores decimales y de miles
	if strings.Contains(amountStr, ",") && strings.Contains(amountStr, ".") {
		// Formato europeo: "1.234,56"
		amountStr = strings.ReplaceAll(amountStr, ".", "")
		amountStr = strings.ReplaceAll(amountStr, ",", ".")
	} else if strings.Count(amountStr, ",") == 1 && !strings.Contains(amountStr, ".") {
		// Ej: "99,99" → formato europeo con coma decimal
		amountStr = strings.ReplaceAll(amountStr, ",", ".")
	} else if strings.Contains(amountStr, ",") {
		// Coma como separador de miles
		amountStr = strings.ReplaceAll(amountStr, ",", "")
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, currency, err
	}

	return amount, currency, nil
}
