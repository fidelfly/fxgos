package cron

import (
	"context"
	"strconv"
	"time"

	"github.com/fidelfly/gox/cronx"
	"github.com/fidelfly/gox/logx"

	"github.com/fidelfly/fxgos/cmd/service/cron/res"
	"github.com/fidelfly/gostool/db"
)

var myCronx = cronx.New(cronx.WithMiddleware(logJob))

var jobMap = make(map[string]cronx.Job)

func Initialize() error {
	err := db.Synchronize(
		new(res.CronJob), new(res.JobAudit),
	)

	//myCronx = cronx.New(cronx.WithMiddleware(logJob))
	return err
}

func Start() {
	jobs := make([]*res.CronJob, 0)
	if err := db.Find(&jobs, db.Where("status = ?", JobEffective)); err != nil {
		if len(jobs) > 0 {
			for _, job := range jobs {
				logx.CaptureError(activateJob(job))
				if job.Type == CyclicJobType && job.OnStartup {
					logx.CaptureError(runJob(job))
				}
			}
		}
	}
	myCronx.Start()
}

func logJob(job cronx.Job) cronx.Job {
	return cronx.FuncJob(func(ctx context.Context) (err error) {
		md := cronx.GetMetadata(ctx)
		id := GetJobId(md)
		code := GetJobCode(md)
		jobType := GetJobType(md)
		audit := &res.JobAudit{
			JobId:     id,
			StartTime: time.Now(),
		}
		logx.Infof(`Job (id = %d, code = "%s", type = "%s") begins to run.`, id, code, jobType)
		err = job.Run(ctx)
		audit.EndTime = time.Now()
		audit.Duration = audit.EndTime.Sub(audit.StartTime)
		if err != nil {
			audit.Status = Failed
			audit.ErrMessage = err.Error()
			logx.Error(`Job (id = %d, code = "%s", type = "%s") ends with error : %s.`, id, code, jobType, err.Error())
		} else {
			audit.Status = Done
			logx.Infof(`Job (id = %d, code = "%s", type = "%s") ends with positive result`, id, code, jobType)
		}
		if jobType == OneTimeJobType {
			if way, ok := md.Get(MetaJobRunWay); ok && way == automatic {
				logx.CaptureError(db.Update(&res.CronJob{
					Id:     id,
					Status: JobExpired,
				}, db.ID(id), db.Cols("status")))
			}
		}
		logx.CaptureError(db.Create(audit))
		return
	})
}

func GetJobId(md *cronx.Metadata) int64 {
	if v, ok := md.Get(MetaJobId); ok {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			return id
		}
	}
	return 0
}

func GetJobCode(md *cronx.Metadata) string {
	if v, ok := md.Get(MetaJobCode); ok {
		return v
	}
	return ""
}

func GetJobType(md *cronx.Metadata) string {
	if v, ok := md.Get(MetaJobType); ok {
		return v
	}
	return ""
}
