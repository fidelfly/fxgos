package system

import (
	"time"
)


type User struct {
	Id int64 `xorm:"autoincr`
	Code string
	Name string
	Password string
	CreateTime time.Time `xorm:"created"`
}

type TraceLog struct {
	Id int64 `xorm:"autoincr"`
	UserId int64
	User string
	RequestUrl string
	LogTime time.Time `xorm:"created"`
	Code string
	Type string
	Message string
	Info string `xorm:"text"`
}

