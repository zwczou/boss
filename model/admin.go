package model

import (
	"time"
)

type Administrator struct {
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
	Remark      string
	DivideModes []DivideMode
	Roles       []Role       `gorm:"many2many:admin_role_admin"`
	Permissions []Permission `gorm:"many2many:admin_admin_permissions"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p Administrator) TableName() string {
	return "administrators"
}

type DivideMode struct {
	Id              int
	AdministratorId int `gorm:"index:"idx_admin_id"`
	Operator        int `gorm:"type:tinyint"`
	Type            int `gorm:"type:tinyint"`
	Value           float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (dm DivideMode) TableName() string {
	return "divide_modes"
}

type Role struct {
	Id             int
	Name           string
	Slug           string
	Administrators []Administrator `gorm:"many2many:admin_role_admins"`
	Permissions    []Permission    `gorm:"many2many:admin_role_permissions"`
	Menus          []Menu          `gorm:"many2many:admin_role_menus"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (r Role) TableName() string {
	return "admin_roles"
}

type Permission struct {
	Id             int
	Name           string `gorm:"size:80"`
	Slug           string `gorm:"size:80;unique_index:uiq_slug"`
	HttpMethod     string
	HttpPath       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Roles          []Role          `gorm:"many2many:admin_role_permissions"`
	Administrators []Administrator `gorm:"many2many:admin_admin_permissions"`
}

func (p Permission) TableName() string {
	return "admin_permissions"
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
