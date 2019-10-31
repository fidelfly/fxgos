package dbo

import (
	"context"

	"github.com/fidelfly/gostool/db"
)

func Delete(ctx context.Context, target interface{}, option []db.StatementOption, hooks ...SessionHook) (int64, error) {
	dbs := CurrentDBSession(ctx, db.AutoClose(true))
	applyHooks(ctx, dbs.Session, hooks...)
	return dbs.Delete(target, option...)
}
