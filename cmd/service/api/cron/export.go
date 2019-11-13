package cron

import (
	"context"
	"fmt"

	"github.com/fidelfly/gox/cronx"

	"github.com/fidelfly/fxgos/cmd/service"
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

func getServer() Service {
	if v, ok := service.GetService(ServiceName); ok {
		if server, ok := v.(Service); ok {
			return server
		}
	}
	panic(fmt.Sprintf("Service(%s) is not registered", ServiceName))
}

func Create(ctx context.Context, opts ...JobOption) (int64, error) {
	return getServer().Create(ctx, opts...)
}
func AddJob(code string, job cronx.Job) {
	getServer().AddJob(code, job)
}
