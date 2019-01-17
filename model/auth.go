package model

import (
	"time"
)

const (
	AuthUnauthed = iota
	AuthAudited
	AuthUnpass
	AuthPass
)

type Auth struct {
	Id         int
	CardId     int   `gorm:"index:idx_card_id"`
	Card       *Card `json:",omitempty"`
	UserId     int   `gorm:"index:idx_user_id"`
	UserIp     string
	Name       string
	IdentNum   string
	Mobile     string
	Cover      string
	Back       string
	People     string
	Status     int `gorm:"type:tinyint"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	AuthTracks []AuthTrack
}

func (a Auth) TableName() string {
	return "auths"
}

type AuthTrack struct {
	Id        int
	AuthId    int `gorm:"index:idx_auth_id"`
	OldStatus int
	Status    int
	Remark    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (at AuthTrack) TableName() string {
	return "auth_tracks"
}
