package model

import (
	"time"
)

type Remark struct {
	Id        int
	CardId    int `gorm:"card_id"`
	OwnerId   int `gorm:"idx_owner_id"`
	OwnerType string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r Remark) TableName() string {
	return "remarks"
}
