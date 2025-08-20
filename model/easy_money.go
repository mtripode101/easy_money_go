package model

import (
	"time"
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
