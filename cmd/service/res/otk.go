package res

import "time"

type OneTimeKey struct {
	Id         int64     `xorm:"autoincr pk id" json:"id"`
	Key        string    `xorm:"not null unique" json:"key"`
	Type       string    `xorm:"not null" json:"type"`
	TypeId     string    `xorm:"not null" json:"type_id"`
	Usage      string    `xorm:"not null" json:"usage"`
	Data       string    `xorm:"not null" json:"data"`
	Consumed   bool      `xorm:"not null" json:"consumed"`
	Invalid    bool      `xorm:"not null" json:"invalid"`
	CreateTime time.Time `xorm:"created" json:"create_time"`
}
