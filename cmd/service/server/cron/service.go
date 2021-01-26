package cron

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/fidelfly/gox/cronx"
	"github.com/fidelfly/gox/errorx"
	"github.com/fidelfly/gox/logx"

	"github.com/fidelfly/fxgos/cmd/service/api/cron"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/gostool/db"
)

type server struct {
}

func (s server) Create(ctx context.Context, opts ...cron.JobOption) (int64, error) {
	cronJob := &res.CronJob{
		Status: cron.JobEffective,
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(cronJob)
		}
	}

	if len(cronJob.Spec) == 0 && cronJob.Timer.IsZero() {
		cronJob.Status = cron.JobInvalid
	}

	if _, err := db.Create(cronJob); err != nil {
		return 0, syserr.DatabaseErr(err)
	}

	if cronJob.Status == cron.JobEffective {
		logx.CaptureError(s.activateJob(cronJob))
	}
	return cronJob.Id, nil
}

func (s server) activateJob(cronJob *res.CronJob) (int, error) {
	if cronJob.Status == cron.JobEffective {
		if id, ok := getRunningJob(cronJob.Id); ok {
			return id, nil
		}
		if job, ok := jobMap[cronJob.Code]; ok {
			meta := s.makeMeta(cronJob)
			meta[cron.MetaJobRunWay] = automatic
			var jobId int
			var err error
			if cronJob.Type == cron.OneTimeJobType {
				if cronJob.Timer.After(time.Now()) {
					jobId = myCronx.AddTimerJob(cronJob.Timer, job, meta)
				} else {
					jobId = myCronx.AddTimerJob(time.Now().Add(2*time.Minute), job, meta)
				}
			} else {
				jobId, err = myCronx.AddJob(cronJob.Spec, job, meta)
			}
			addRunningJob(cronJob.Id, jobId)
			return jobId, err
		} else {
			return 0, errorx.NewError(
				"err.cron.job_definition_not_found",
				fmt.Sprintf("can't find job definition with code = %s", cronJob.Code),
			)
		}
	} else {
		return 0, errorx.NewError(
			"err.cron.job_is_invalid",
			fmt.Sprintf("job with code = %s is invalid", cronJob.Code),
		)
	}
}

func (s server) runJob(cronJob *res.CronJob) error {
	if job, ok := jobMap[cronJob.Code]; ok {
		meta := s.makeMeta(cronJob)
		meta[cron.MetaJobRunWay] = manual
		return myCronx.RunJob(job, meta)
	} else {
		return errorx.NewError(
			"err.cron.job_definition_not_found",
			fmt.Sprintf("can't find job definition with code = %s", cronJob.Code),
		)
	}
}

const (
	automatic = "automatic"
	manual    = "manual"
)

func (s server) makeMeta(cronJob *res.CronJob) map[string]string {
	m := make(map[string]string)
	m[cron.MetaJobId] = strconv.FormatInt(cronJob.Id, 10)
	m[cron.MetaJobCode] = cronJob.Code
	m[cron.MetaJobType] = cronJob.Type
	m[cron.MetaJobTraceable] = strconv.FormatBool(cronJob.Traceable)
	if len(cronJob.Meta) > 0 {
		for k, v := range cronJob.Meta {
			m[k] = v
		}
	}
	return m
}

func (s server) AddJob(code string, job cronx.Job) {
	if _, ok := jobMap[code]; ok {
		panic(fmt.Sprintf("duplicated job code : %s", code))
	}
	jobMap[code] = job
}
