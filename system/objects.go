package system

import (
	"time"
)


type User struct {
	Id int64 `xorm:"autoincr`
	Code string
	Name string
	Avatar int64
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

type Assets struct {
	Id int64 `xorm:"autoincr"`
	Md5 string
	Type string
	Size int64
	Name string
	Data []byte
	CreateTime time.Time `xorm:"created"`
}

