package audit

import (
	"context"

	"github.com/fidelfly/fxgos/cmd/service"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/gostool/dbo"
)

const (
	ServiceName = "service.audit"
)

func RegisterServer(server Service) {
	service.Register(ServiceName, server)
}

type Service interface {
	ListTrail(ctx context.Context, input *dbo.ListInfo, conds ...string) ([]*res.Systrail, int64, error)
}
