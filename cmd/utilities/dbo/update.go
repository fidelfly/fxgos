package dbo

import (
	"context"

	"github.com/fidelfly/gostool/db"
)

type UpdateOption func(ctx context.Context, target interface{}) []db.QueryOption

func Update(ctx context.Context, target interface{}, hooks ...SessionHook) error {
}
