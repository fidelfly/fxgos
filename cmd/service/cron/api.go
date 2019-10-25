package cron

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/fidelfly/gox/cronx"
	"github.com/fidelfly/gox/errorx"
	"github.com/fidelfly/gox/logx"

	"github.com/fidelfly/fxgos/cmd/service/cron/res"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/gostool/db"
)

const (
	OneTimeJobType = "cron.one-time-job"
	CyclicJobType  = "cron.cyclic-job"
)

const (
	JobInvalid = iota
	JobEffective
	JobExpired
)

const (
	Failed = iota
	Done
)

type JobOption func(*res.CronJob)

func JobCode(code string) JobOption {
	return func(job *res.CronJob) {
		job.Code = code
	}
}
func CyclicJob(spec string) JobOption {
	return func(job *res.CronJob) {
		job.Spec = spec
		job.Type = CyclicJobType
	}
}
func OneTimeJob(t time.Time) JobOption {
	return func(job *res.CronJob) {
		job.Timer = t
		job.Type = OneTimeJobType
	}
}
func JobType(t string) JobOption {
	return func(job *res.CronJob) {
		job.Type = t
	}
}
func JobMeta(data ...map[string]string) JobOption {
	return func(job *res.CronJob) {
		md := make(map[string]string)
		for _, dm := range data {
			if len(dm) > 0 {
				for k, v := range dm {
					md[k] = v
				}
			}
		}
		job.Meta = md
	}
}

func StartupJob() JobOption {
	return func(job *res.CronJob) {
		job.OnStartup = true
	}
}

func Create(ctx context.Context, opts ...JobOption) (int64, error) {
	cronJob := &res.CronJob{
		Status: JobEffective,
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(cronJob)
		}
	}

	if len(cronJob.Spec) == 0 && cronJob.Timer.IsZero() {
		cronJob.Status = JobInvalid
	}

	if _, err := db.Create(cronJob); err != nil {
		return 0, syserr.DatabaseErr(err)
	}

	if cronJob.Status == JobEffective {
		logx.CaptureError(activateJob(cronJob))
	}
	return cronJob.Id, nil
}

func activateJob(cronJob *res.CronJob) (int, error) {
	if cronJob.Status == JobEffective {
		if job, ok := jobMap[cronJob.Code]; ok {
			meta := makeMeta(cronJob)
			meta[MetaJobRunWay] = automatic
			if cronJob.Type == OneTimeJobType {
				return myCronx.AddTimerJob(cronJob.Timer, job, meta), nil
			} else {
				return myCronx.AddJob(cronJob.Spec, job, meta)
			}
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

func runJob(cronJob *res.CronJob) error {
	if job, ok := jobMap[cronJob.Code]; ok {
		meta := makeMeta(cronJob)
		meta[MetaJobRunWay] = manual
		return myCronx.RunJob(job, meta)
	} else {
		return errorx.NewError(
			"err.cron.job_definition_not_found",
			fmt.Sprintf("can't find job definition with code = %s", cronJob.Code),
		)
	}
}

const (
	MetaJobId     = "meta.job.id"
	MetaJobCode   = "meta.job.code"
	MetaJobType   = "meta.job.type"
	MetaJobRunWay = "meta.job.run.way"
)

const (
	automatic = "automatic"
	manual    = "manual"
)

func makeMeta(cronJob *res.CronJob) map[string]string {
	m := make(map[string]string)
	m[MetaJobId] = strconv.FormatInt(cronJob.Id, 10)
	m[MetaJobCode] = cronJob.Code
	m[MetaJobType] = cronJob.Type
	if len(cronJob.Meta) > 0 {
		for k, v := range cronJob.Meta {
			m[k] = v
		}
	}
	return m
}

func AddJob(code string, job cronx.Job) {
	if _, ok := jobMap[code]; ok {
		panic(fmt.Sprintf("duplicated job code : %s", code))
	}
	jobMap[code] = job
}
