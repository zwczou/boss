package model

import (
	"time"
)

type Administrator struct {
	Id          int
	ParentId    int
	Username    string `gorm:"unique_index:uiq_username"`
	Mobile      string `gorm:"unique_index:uiq_mobile"`
	Name        string
	Password    string
	IsDisabled  bool
	Roles       []Role       `gorm:"many2many:admin_role_users"`
	Permissions []Permission `gorm:"many2many:admin_user_permissions"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (a Administrator) TableName() string {
	return "admin_users"
}

func (a Administrator) CheckPassword(pwd string) bool {
	if len(a.Password) == 40 {
		return a.Password == encryptPassword(pwd)
	}
	return comparePasswords(a.Password, []byte(pwd))
}

func (a *Administrator) SetPassword(pwd string) {
	a.Password = hashAndSalt([]byte(pwd))
}

type Role struct {
	Id             int
	Name           string
	Slug           string
	Administrators []Administrator `gorm:"many2many:admin_role_users"`
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
	Administrators []Administrator `gorm:"many2many:admin_user_permissions"`
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
