package system

import (
	"time"
)

type Tenant struct {
	Id int64 `xorm:"autoincr"`
	Code string
	Name string
	CreateTime time.Time `xorm:"created"`
}

type User struct {
	Id int64 `xorm:"autoincr`
	Code string
	Name string
	Password string
	TenantId int64
	CreateTime time.Time `xorm:"created"`
}

type TraceLog struct {
	Id int64 `xorm:"autoincr"`
	UserId int64
	User string
	TenantId int64
	Tenant string
	RequestUrl string
	LogTime time.Time `xorm:"created"`
	Code string
	Type string
	Message string
	Info string `xorm:"text"`
}

//func init() {
//	database.AttachTable(new(Tenant), new(User), new(TraceLog))
//}
