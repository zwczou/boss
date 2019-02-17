package model

import (
	"time"
)

type Partner struct {
	Id          int
	ParentId    int
	Nickname    string `gorm:"unique_index:uiq_nickname"`
	Mobile      string `gorm:"unique_index:uiq_mobile"`
	Password    string
	Company     string
	Email       string
	Contact     string
	Address     string
	IsDisabled  bool
	Level       int
	CardCount   int
	Note        string
	DivideModes []DivideMode
	Roles       []Role       `gorm:"many2many:admin_role_users"`
	Permissions []Permission `gorm:"many2many:admin_user_permissions"`
	Cards       []Card       `gorm:"many2many:partner_cards"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p Partner) TableName() string {
	return "admin_users"
}

func (p Partner) CheckPassword(pwd string) bool {
	if len(p.Password) == 40 {
		return p.Password == encryptPassword(pwd)
	}
	return ComparePasswords(p.Password, []byte(pwd))
}

func (p *Partner) SetPassword(pwd string) {
	p.Password = encryptPassword(pwd)
}

type Administrator = Partner

type Role struct {
	Id          int
	Name        string
	Slug        string
	Partners    []Partner    `gorm:"many2many:admin_role_users"`
	Permissions []Permission `gorm:"many2many:admin_role_permissions"`
	Menus       []Menu       `gorm:"many2many:admin_role_menus"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r Role) TableName() string {
	return "admin_roles"
}

type Permission struct {
	Id         int
	Name       string `gorm:"size:80"`
	Slug       string `gorm:"size:80;unique_index:uiq_slug"`
	HttpMethod string
	HttpPath   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Roles      []Role    `gorm:"many2many:admin_role_permissions"`
	Partners   []Partner `gorm:"many2many:admin_user_permissions"`
}

func (p Permission) TableName() string {
	return "user_permissions"
}

type Menu struct {
	Id        int
	ParentId  int
	Order     int
	Title     string `gorm:"size:50"`
	Icon      string `gorm:"size:50"`
	Uri       string `gorm:"size:50"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Roles     []Role  `gorm:"many2many:admin_role_menus"`
	Menus     []*Menu `gorm:"-"`
}

func (m Menu) TableName() string {
	return "admin_menus"
}
