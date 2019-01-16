package model

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	OrderWeixin = iota
	OrderAlipay
)

const (
	OrderCreated = iota
	OrderPaid
	OrderClosed
)

const (
	OrderPlan = iota
	OrderBatch
)

const (
	SourceWeixinWeb = iota + 1
	SourceWebBatch
)

type Order struct {
	Id            int
	PartnerId     int // 只是用来查询是那个微信支付信息
	No            string
	PlanId        int          `gorm:"index:idx_plan_id"`
	Plan          *Plan        `json:",omitempty"`
	CardId        int          `gorm:"index:idx_card_id"`
	Card          *Card        `json:",omitempty"`
	PartnerPlanId int          `gorm:"index:idx_partner_plan_id"`
	PartnerPlan   *PartnerPlan `json:",omitempty"`
	Type          int
	Source        int
	UserIp        string
	Name          string
	Data          float64
	MonthData     float64
	Voice         int
	MonthVoice    int
	Day           int
	Price         float64
	IsNext        bool
	Status        int
	PayType       int
	CreatedAt     time.Time
	PaidAt        *time.Time
	ClosedAt      *time.Time
	UpdatedAt     time.Time
}

func (o Order) TableName() string {
	return "orders"
}

// 创建订单的时候自动填充订单号
func (o *Order) BeforeSave(scope *gorm.Scope) {
	if noField, ok := scope.FieldByName("No"); ok {
		if noField.IsBlank {
			no := fmt.Sprintf("T%s%04d", time.Now().Format("20060102150405"), 1000+rand.Intn(8999))
			noField.Set(no)
		}
	}
}

func (o *Order) BeforeCreate(scope *gorm.Scope) {
	o.BeforeSave(scope)
}

type Comission struct {
	Id        int
	PartnerId int      `gorm:"unique_index:uiq_partner_order"`
	Partner   *Partner `json:",omitempty"`
	OrderId   int      `gorm:"unique_index:uiq_partner_order"`
	Order     *Order   `json:",omitempty"`
	Amount    float64
	NetAmount float64
	CreatedAt time.Time
}

func (c Comission) TableName() string {
	return "comissions"
}
