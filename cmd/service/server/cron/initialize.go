package cron

import (
	"context"
	"sync"
	"time"

	"github.com/fidelfly/gox/cronx"
	"github.com/fidelfly/gox/logx"

	"github.com/fidelfly/fxgos/cmd/service/api/cron"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/gostool/db"
	"github.com/fidelfly/gostool/dbo"
)

var myCronx = cronx.New(cronx.WithMiddleware(logJob))
var s *server
var jobMap = make(map[string]cronx.Job)
var runningJobs = make(map[int64]int)
var runningLock = sync.RWMutex{}

func Initialize() error {
	err := db.Synchronize(
		new(res.CronJob), new(res.JobAudit),
	)
	s = &server{}
	cron.RegisterServer(s)
	return err
}

func Start() {
	jobs := make([]*res.CronJob, 0)
	if err := db.Find(&jobs, db.Where("status = ?", cron.JobEffective)); err != nil {
		if len(jobs) > 0 {
			for _, job := range jobs {
				logx.CaptureError(s.activateJob(job))
				if job.Type == cron.CyclicJobType && job.OnStartup {
					logx.CaptureError(s.runJob(job))
				}
			}
		}
	}
	myCronx.Start()
}

func logJob(job cronx.Job) cronx.Job {
	return cronx.FuncJob(func(ctx context.Context) (err error) {
		md := cronx.GetMetadata(ctx)
		id := cron.GetJobId(md)
		code := cron.GetJobCode(md)
		jobType := cron.GetJobType(md)
		audit := &res.JobAudit{
			JobId:     id,
			StartTime: time.Now(),
		}
		logx.Infof(`Job (id = %d, code = "%s", type = "%s") begins to run.`, id, code, jobType)
		err = job.Run(ctx)
		audit.EndTime = time.Now()
		audit.Duration = audit.EndTime.Sub(audit.StartTime)
		if err != nil {
			audit.Status = cron.Failed
			audit.ErrMessage = err.Error()
			logx.Error(`Job (id = %d, code = "%s", type = "%s") ends with error : %s.`, id, code, jobType, err.Error())
		} else {
			audit.Status = cron.Done
			logx.Infof(`Job (id = %d, code = "%s", type = "%s") ends with positive result`, id, code, jobType)
		}
		dbs := dbo.CurrentDBSession(ctx, dbo.DefaultSession)
		defer dbs.Close()
		if jobType == cron.OneTimeJobType {
			if way, ok := md.Get(cron.MetaJobRunWay); ok && way == automatic {
				logx.CaptureError(dbs.Update(&res.CronJob{
					Id:     id,
					Status: cron.JobExpired,
				}, db.ID(id), db.Cols("status")))
			}
			removeRunningJob(id)
		}
		logx.CaptureError(dbs.Insert(audit))
		return
	})
}

func addRunningJob(resId int64, jobId int) {
	runningLock.Lock()
	defer runningLock.Unlock()

	runningJobs[resId] = jobId
}

func removeRunningJob(resId int64) {
	runningLock.Lock()
	defer runningLock.Unlock()

	delete(runningJobs, resId)
}

func getRunningJob(resId int64) (int, bool) {
	runningLock.RLock()
	defer runningLock.RUnlock()

	if id, ok := runningJobs[resId]; ok {
		return id, true
	} else {
		return 0, false
	}
}
