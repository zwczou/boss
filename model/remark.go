package model

import (
	"time"
)

type Remark struct {
	Id        int
	CardId    int    `gorm:"unique_index:idx_card_owner"`
	OwnerId   int    `gorm:"unique_index:idx_card_owner"`
	OwnerType string `gorm:"unique_index:idx_card_owner"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r Remark) TableName() string {
	return "remarks"
}
