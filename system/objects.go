package system

import (
	"time"
)

type User struct {
	ID         int64 `xorm:"autoincr id"`
	Code       string
	Name       string
	Avatar     int64
	Password   string
	CreateTime time.Time `xorm:"created"`
}

type TraceLog struct {
	ID         int64 `xorm:"autoincr id"`
	UserID     int64
	User       string    `xorm:"user_id"`
	RequestURL string    `xorm:"request_url"`
	LogTime    time.Time `xorm:"created"`
	Code       string
	Type       string
	Message    string
	Info       string `xorm:"text"`
}

type Assets struct {
	ID         int64 `xorm:"autoincr id"`
	Md5        string
	Type       string
	Size       int64
	Name       string
	Data       []byte
	CreateTime time.Time `xorm:"created"`
}
