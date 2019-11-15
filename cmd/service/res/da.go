package res

import "time"

type SecurityGroup struct {
	Id          int64     `xorm:"autoincr pk id" json:"id"`
	Code        string    `xorm:"not null unique" json:"code"`
	Description string    `json:"description"`
	CreateTime  time.Time `xorm:"created" json:"-"`
	UpdateTime  time.Time `xorm:"updated" json:"-"`
}

type UserSg struct {
	Id            int64 `xorm:"autoincr pk id" json:"id"`
	UserId        int64 `json:"user_id"`
	SecurityGroup int64 `json:"security_group"`
}

type ResourceSg struct {
	Id            int64  `xorm:"autoincr pk id" json:"id"`
	ResType       string `json:"res_type"`
	ResId         int64  `json:"res_id"`
	SecurityGroup int64  `json:"security_group"`
}
