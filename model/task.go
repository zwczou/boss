package model

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	TaskPending = iota
	TaskStarted
	TaskDone
)

const (
	TypeLoadCard = iota
	TypeDistCard
	TypeBatchCharge
	TypeLoadDevice
	TypeExportCard
	TypeExportOrder
	TypeExportDevice
	TypeExportComission
)

type Task struct {
	Id        int
	No        string   `gorm:"size:25"`
	PartnerId int      `gorm:"index:idx_partner_id"`
	Partner   *Partner `json:",omitempty"`
	Type      int
	Status    int
	Spent     int
	Total     int
	Success   int
	Remark    string
	StartedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t Task) TableName() string {
	return "tasks"
}

func (t *Task) BeforeSave(scope *gorm.Scope) {
	if noField, ok := scope.FieldByName("No"); ok {
		if noField.IsBlank {
			no := fmt.Sprintf("T%s%04d", time.Now().Format("20060102150405"), 1000+rand.Intn(8999))
			noField.Set(no)
		}
	}
}

func (t *Task) BeforeCreate(scope *gorm.Scope) {
	t.BeforeSave(scope)
}

func (t *Task) Started() {
	var now = time.Now()
	t.Status = TaskStarted
	t.StartedAt = &now
}

func (t *Task) Done() {
	t.Status = TaskDone
}
