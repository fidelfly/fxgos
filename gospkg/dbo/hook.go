package dbo

import (
	"context"

	"github.com/fidelfly/gostool/db"
)

type StatementHook func(ctx context.Context, bean interface{})

func StatementBeforeHook(ctx context.Context, hook StatementHook) db.StatementOption {
	return db.BeforeClosure(func(target interface{}) {
		hook(ctx, target)
	})
}

func StatementAfterHook(ctx context.Context, hook StatementHook) db.StatementOption {
	return db.AfterClosure(func(target interface{}) {
		hook(ctx, target)
	})
}
