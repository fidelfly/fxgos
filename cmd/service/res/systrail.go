package res

import "time"

type Systrail struct {
	Id                int64             `xorm:"autoincr pk" json:"id"`
	Key               string            `json:"-"`
	Code              string            `json:"code"`
	Operation         string            `json:"operation"`
	StartTime         time.Time         `json:"-"`
	EndTime           time.Time         `json:"end_time"`
	Duration          int64             `json:"-"`
	ExecUser          int64             `json:"exec_user"`
	Status            int64             `json:"-"`
	StatusDescription string            `json:"-"`
	Info              map[string]string `xorm:"json" json:"info"`
	RequestId         string            `json:"request_id"`
}
