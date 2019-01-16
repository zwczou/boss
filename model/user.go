package model

import (
	"time"
)

type User struct {
	Id        int
	PartnerId int      `gorm:"index:idx_partner_id"`
	Partner   *Partner `json:",omitempty"`
	WeixinId  string
	AvatarRaw string
	Avatar    string
	Gender    int
	Nickname  *string `gorm:"unique_index:uiq_nickname"`
	Mobile    *string `gorm:"unique_index:uiq_mobile"`
	Password  string
	Cards     []Card `gorm:"many2many:user_cards"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) TableName() string {
	return "users"
}
