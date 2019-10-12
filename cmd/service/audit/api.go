package audit

import (
	"context"

	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/service/audit/res"
)

func ListTrail(ctx context.Context, input db.ListInfo) ([]*res.Systrail, int64, error) {
	resSystrails := make([]*res.Systrail, 0)
	opts := make([]db.QueryOption, 0)
	if len(input.Cond) > 0 {
		opts = append(opts, db.Where(input.Cond))
	}
	queOpts := append(append(db.GetPagingOption(input), db.Desc("end_time")), opts...)
	if err := db.Find(&resSystrails, queOpts...); err != nil {
		return nil, 0, err
	}

	count := int64(len(resSystrails))
	if !(count < input.Results && input.Page == 1) {
		count, _ = db.Count(new(res.Systrail), opts...)
	}
	return resSystrails, count, nil
}
