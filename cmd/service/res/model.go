package res

import "time"

type Model struct {
	Id           int64 `xorm:"autoincr pk id"`
	ResourceType string
	Data         []byte    `xorm:"longblob"`
	Policy       []byte    `xorm:"longblob"`
	CreateTime   time.Time `xorm:"created"`
	UpdateTime   time.Time `xorm:"updated"`
}

type Policy struct {
	Id           int64     `xorm:"autoincr pk id" json:"id"`
	RoleId       int64     `json:"role_id"`
	ResourceType string    `json:"resource_type"`
	Sub          string    `json:"sub"`
	Obj          string    `json:"obj"`
	Act          []string  `xorm:"json" json:"act"`
	CreateTime   time.Time `xorm:"created" json:"-"`
	UpdateTime   time.Time `xorm:"updated" json:"-"`
}
