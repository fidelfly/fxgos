package dbo

import (
	"context"

	"github.com/fidelfly/gostool/db"
)

func Create(ctx context.Context, target interface{}, options ...db.StatementOption) error {
	dbs := CurrentDBSession(ctx, db.AutoClose(true))
	if _, err := dbs.Insert(target, options...); err != nil {
		return err
	}
	return nil
}
