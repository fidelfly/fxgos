package res

import "time"

type User struct {
	Id             int64     `xorm:"autoincr pk id" json:"id"`
	Code           string    `xorm:"not null unique" json:"code"`
	Name           string    `xorm:"not null" json:"name"`
	Email          string    `xorm:"not null unique" json:"email,omitempty"`
	Password       string    `xorm:"not null" json:"password,omitempty"`
	Avatar         int64     `xorm:"not null" json:"avatar"`
	Roles          []int64   `xorm:"json" json:"roles"`
	SecurityGroups []int64   `xorm:"json" json:"security_groups"`
	Region         string    `xorm:"not null" json:"region"`
	Dept           string    `xorm:"not null" json:"dept"`
	Tel            string    `xorm:"not null" json:"tel"`
	Status         int64     `xorm:"not null" json:"status"`
	SuperAdmin     bool      `json:"sa"`
	CreateTime     time.Time `xorm:"created" json:"-"`
	UpdateTime     time.Time `xorm:"updated" json:"-"`
}
