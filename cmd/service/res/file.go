package res

import "time"

type File struct {
	Id         int64  `xorm:"autoincr pk id"`
	Name       string `xorm:"not null"`
	Md5        string `xorm:"not null"`
	Data       []byte `xorm:"longblob not null"`
	Type       string
	Size       int64
	CreateTime time.Time `xorm:"created"`
	CreateUser int64
}
