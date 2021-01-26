package res

import "time"

type CronJob struct {
	Id         int64             `xorm:"autoincr pk id"`
	Code       string            // job classification
	Type       string            // one-time-job, cyclic-job
	Spec       string            // cron specification
	Timer      time.Time         // timer for one-time-job
	Meta       map[string]string `xorm:"json"` // json format of any interface which is the only parameter of the cron job function
	OnStartup  bool              // indicate whether this job will run when system start.
	CreateTime time.Time         `xorm:"created"`
	Status     int               // cron.JobExpired cron.JobEffective cron.JobInvalid
	Traceable  bool              `xorm:"default 1"`
}

type JobAudit struct {
	Id         int64 `xorm:"autoincr pk id"`
	JobId      int64
	StartTime  time.Time
	EndTime    time.Time
	Duration   time.Duration
	Status     int    //cron.Failed cron.Done
	ErrMessage string `xorm:"varchar(1024)"`
}
