package dbo

import (
	"context"

	"github.com/fidelfly/gostool/db"
)

type StatementHook func(ctx context.Context, bean interface{})

func StatementBeforeHook(ctx context.Context, hook StatementHook) db.StatementOption {
	return func(session *db.Session) {
		session.GetXorm().Before(func(target interface{}) {
			hook(ctx, target)
		})
	}
}

func StatementAfterHook(ctx context.Context, hook StatementHook) db.StatementOption {
	return func(session *db.Session) {
		session.GetXorm().Before(func(target interface{}) {
			hook(ctx, target)
		})
	}
}
