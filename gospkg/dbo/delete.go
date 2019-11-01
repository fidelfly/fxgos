package dbo

import (
	"context"

	"github.com/fidelfly/gostool/db"
)

func Delete(ctx context.Context, target interface{}, option ...db.StatementOption) (int64, error) {
	dbs := CurrentDBSession(ctx, db.AutoClose(true))
	return dbs.Delete(target, option...)
}
