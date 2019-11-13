package cron

import (
	"context"

	"github.com/fidelfly/gox/cronx"

	"github.com/fidelfly/fxgos/cmd/service"
)

const (
	ServiceName = "service.cron"
)

func RegisterServer(server Service) {
	service.Register(ServiceName, server)
}

type Service interface {
	Create(ctx context.Context, opts ...JobOption) (int64, error)
	AddJob(code string, job cronx.Job)
}
