package model

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
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

func (p Partner) encryptPassword(pwd string) string {
	key := []byte("zouweicheng@gmail.com")
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(pwd))
	h := sha1.New()
	io.WriteString(h, hex.EncodeToString(mac.Sum(nil)))
	return hex.EncodeToString(h.Sum(nil))
}

func (p Partner) CheckPassword(pwd string) bool {
	if len(p.Password) == 40 {
		return p.Password == p.encryptPassword(pwd)
	}
	return ComparePasswords(p.Password, []byte(pwd))
}

func HashAndSalt(pwd []byte) string {
	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.WithError(err).Error("bcrypt.GenerateFromPassword")
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.WithError(err).Error("bcrypt.CompareHashAndPassword")
		return false
	}
	return true
}

func (p *Partner) SetPassword(pwd string) {
	p.Password = p.encryptPassword(pwd)
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
