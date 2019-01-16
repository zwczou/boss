package model

import (
	"time"
)

type DivideMode struct {
	Id        int
	PartnerId int `gorm:"index:idx_partner_id"`
	Operator  int `gorm:"type:tinyint"`
	Type      int `gorm:"type:tinyint"`
	Value     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (dm DivideMode) TableName() string {
	return "divide_modes"
}
