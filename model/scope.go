package model

import (
	"github.com/jinzhu/gorm"
)

func QueryAdministratorScope(db *gorm.DB) *gorm.DB {
	return db.Preload("Roles.Permissions").Preload("Roles.Menus").Preload("Permissions")
}
