package model

import (
	"github.com/jinzhu/gorm"
)

func QueryAdministratorScope(db *gorm.DB) *gorm.DB {
	return db.Preload("Roles.Permissions").Preload("Roles.Menus").Preload("Permissions")
}

func QueryRoleScope(db *gorm.DB) *gorm.DB {
	return db.Preload("Permissions").Preload("Administrators").Preload("Menus")
}

func QueryPermissionScope(db *gorm.DB) *gorm.DB {
	return db.Preload("Roles").Preload("Administrators")
}

func QueryMenuScope(db *gorm.DB) *gorm.DB {
	return db.Preload("Roles")
}
