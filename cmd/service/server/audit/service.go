package audit

import (
	"context"

	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/gostool/db"
	"github.com/fidelfly/gostool/dbo"
)

type server struct {
}

func (s server) ListTrail(ctx context.Context, input *dbo.ListInfo, conds ...string) ([]*res.Systrail, int64, error) {
	resSystrails := make([]*res.Systrail, 0)

	count, err := dbo.List(ctx, &resSystrails, input, db.Condition(conds...), db.Desc("end_time"))
	if err != nil {
		return nil, 0, err
	}
	return resSystrails, count, nil
}
