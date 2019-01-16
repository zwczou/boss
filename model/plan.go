package model

import (
	"time"
)

const (
	PlanKeep = iota
	PlanReset
	PlanGo
)

const (
	PlanWait = iota
	PlanUsed
	PlanDone
)

type Plan struct {
	Id          int
	Source      int
	Type        int
	Name        string
	Data        float64
	MonthData   float64
	Voice       int
	MonthVoice  int
	Day         int
	Price       float64
	MarketPrice float64
	IntoPrice   float64
	ThirdId     string
	Description string
	IsUnlimited bool
	IsRecommend bool
	IsDeleted   bool
	IsBase      bool // 是否是基础套餐，移动专用
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p Plan) TableName() string {
	return "plans"
}

type PartnerPlan struct {
	Id          int
	PlanId      int      `gorm:"index:idx_plan_id"`
	Plan        *Plan    `json:",omitempty"`
	PartnerId   int      `gorm:"index:idx_partner_id"`
	Partner     *Partner `json:",omitempty"`
	Source      int
	Price       float64
	MarketPrice float64
	IntoPrice   float64
	Remark      string
	Priority    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (pp PartnerPlan) TableName() string {
	return "partner_plans"
}

type CardPlan struct {
	Id             int
	CardId         int   `gorm:"index:idx_card_id"`
	Card           *Card `json:",omitempty"`
	PlanId         int   `gorm:"index:idx_plan_id"`
	Plan           *Plan `json:",omitmempty"`
	OrderId        int
	Used           float64
	TotalUsed      float64
	Total          float64
	VoiceUsed      int
	VoiceTotalUsed int
	VoiceTotal     int
	Month          int
	Status         int
	ThirdId        string
	ThirdNo        string
	ActivatedAt    *time.Time
	ExpiredAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (c CardPlan) TableName() string {
	return "card_plans"
}
