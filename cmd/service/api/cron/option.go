package cron

import (
	"time"

	"github.com/fidelfly/fxgos/cmd/service/res"
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

func RunOnStartup() JobOption {
	return func(job *res.CronJob) {
		job.OnStartup = true
	}
}
