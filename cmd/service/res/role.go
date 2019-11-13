package res

import "time"

type Role struct {
	Id          int64     `xorm:"autoincr pk id" json:"id"`
	Code        string    `xorm:"not null unique" json:"code"`
	Roles       []int64   `xorm:"json" json:"roles"`
	Description string    `json:"description"`
	CreateTime  time.Time `xorm:"created" json:"-"`
	UpdateTime  time.Time `xorm:"updated" json:"-"`
}
