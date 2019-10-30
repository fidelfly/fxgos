package dbo

import (
	"context"

	"github.com/fidelfly/fxgos/cmd/utilities/mctx"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/gostool/db"
)

func Create(ctx context.Context, target interface{}, hooks ...SessionHook) error {
	if target == nil {
		return syserr.ErrInvalidParam
	}

	dbs := mctx.CurrentDBSession(ctx, db.AutoClose(true))
	applyHooks(ctx, dbs, hooks...)
	if _, err := dbs.Insert(target); err != nil {
		return syserr.DatabaseErr(err)
	}
}
