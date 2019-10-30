package dbo

import (
	"context"

	"github.com/fidelfly/gostool/db"
)

func Create(ctx context.Context, target interface{}, hooks ...SessionHook) error {
	dbs := CurrentDBSession(ctx, db.AutoClose(true))
	applyHooks(ctx, dbs.Session, hooks...)
	if _, err := dbs.Insert(target); err != nil {
		return err
	}
	return nil
}
