package model

import (
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Administrator struct {
	Id          int
	ParentId    int
	Username    string `gorm:"unique_index:uiq_username"`
	Mobile      string `gorm:"unique_index:uiq_mobile"`
	Name        string
	Password    string
	Level       int
	IsDisabled  bool
	Roles       []Role       `gorm:"many2many:admin_role_users"`
	Permissions []Permission `gorm:"many2many:admin_user_permissions"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (a Administrator) TableName() string {
	return "admin_users"
}

// 创建密码
func (a Administrator) CheckPassword(pwd string) bool {
	if len(a.Password) == 40 {
		return a.Password == encryptPassword(pwd)
	}
	return comparePasswords(a.Password, []byte(pwd))
}

// 设置密码
func (a *Administrator) SetPassword(pwd string) {
	a.Password = hashAndSalt([]byte(pwd))
}

// 是否为超级管理员
func (a *Administrator) IsSuperUser() bool {
	if a.Id == 1 {
		return true
	}
	for _, role := range a.Roles {
		if role.Id == 1 {
			return true
		}
	}
	return false
}

// 获取所有权限
func (a *Administrator) Permission() (permissions []Permission) {
	permissions = append(permissions, a.Permissions...)
	var exists = make(map[int]bool)
	for _, role := range a.Roles {
		for _, m := range role.Permissions {
			if _, ok := exists[m.Id]; !ok {
				permissions = append(permissions, m)
				exists[m.Id] = true
			}
		}
	}
	return
}

// 检测请求权限
// 先检测请求方法
// 然后检测请求路径
func (a *Administrator) Check(req *http.Request) bool {
	if a.IsSuperUser() {
		return true
	}

	r := strings.NewReplacer("*", ".*", "?", ".?")
	for _, perm := range a.Permission() {
		if perm.HttpMethod == "" || strings.Contains(perm.HttpMethod, req.Method) {
			paths := strings.Split(perm.HttpPath, "\n")
			for _, path := range paths {
				path = r.Replace(strings.TrimSpace(path))
				if regexp.MustCompile(path).FindString(req.URL.Path) == req.URL.Path {
					return true
				}
			}
		}
	}
	return false
}

// 获取管理员菜单
func (a *Administrator) Menu() (menus Menus) {
	if a.IsSuperUser() {
		db.Find(&menus)
		return
	}
	var exists = make(map[int]bool)
	for _, role := range a.Roles {
		for _, m := range role.Menus {
			if _, ok := exists[m.Id]; !ok {
				menus = append(menus, m)
				exists[m.Id] = true
			}
		}
	}
	return
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

func (m Menu) InMenu(path string) bool {
	if m.Uri == path {
		return true
	}
	for _, ms := range m.Menus {
		if ok := ms.InMenu(path); ok {
			return true
		}
	}
	return false
}

func (m *Menu) Sort() {
	if len(m.Menus) <= 1 {
		return
	}
	sort.SliceStable(m.Menus, func(i, j int) bool {
		if m.Menus[i].Order == m.Menus[j].Order {
			return m.Menus[i].Id < m.Menus[j].Id
		}
		return m.Menus[i].Order < m.Menus[j].Order
	})
}

type Menus []Menu

// 将菜单转换为Tree模式
func (ms Menus) Tree() Menus {
	var tempMenus = make(map[int]*Menu)
	for idx := range ms {
		m := ms[idx]
		tempMenus[m.Id] = &m
	}
	return menuTree(tempMenus)
}

// 转换菜单为tree形式
func menuTree(ms map[int]*Menu) (menus Menus) {
	for _, m := range ms {
		if val, ok := ms[m.ParentId]; ok {
			val.Menus = append(val.Menus, m)
		}
	}
	for _, m := range ms {
		m.Sort()
		if m.ParentId == 0 {
			menus = append(menus, *m)
		}
	}
	sort.SliceStable(menus, func(i, j int) bool {
		if menus[i].Order == menus[j].Order {
			return menus[i].Id < menus[j].Id
		}
		return menus[i].Order < menus[j].Order
	})
	return
}
