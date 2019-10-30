package mctx

import (
	"context"

	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/gostool/db"
)

type dbSessionKey struct{}

type CtxSession struct {
	dbs        *db.Session
	controlled bool
}

func (cs *CtxSession) Begin() error {
	if cs.controlled {
		if err := cs.dbs.BeginTransaction(); err != nil {
			return syserr.DatabaseErr(err)
		}
	}
	return nil
}

func (cs *CtxSession) Commit() error {
	if cs.controlled {
		if err := cs.dbs.Commit(); err != nil {
			return syserr.DatabaseErr(err)
		}
	}
	return nil
}

func (cs *CtxSession) Rollback() error {
	if cs.controlled {
		if err := cs.dbs.Rollback(); err != nil {
			return syserr.DatabaseErr(err)
		}
	}
	return nil
}

func (cs *CtxSession) Close() {
	if cs.controlled {
		cs.dbs.Close()
	}
}

/*func WithDBSession(ctx context.Context, dbs *db.Session) context.Context {
	return context.WithValue(ctx, dbSessionKey{}, dbs)
}*/

func WithDBSession(ctx context.Context, opts ...db.SessionOption) (context.Context, *CtxSession) {
	dbs := CurrentDBSession(ctx)
	if dbs != nil {
		return ctx, &CtxSession{dbs, false}
	}
	dbs = db.NewSession(opts...)
	return context.WithValue(ctx, dbSessionKey{}, dbs), &CtxSession{dbs, true}
}

func CurrentDBSession(ctx context.Context, opts ...db.SessionOption) *db.Session {
	if v := ctx.Value(dbSessionKey{}); v != nil {
		if dbs, ok := v.(*db.Session); ok {
			return dbs
		}
	}
	if len(opts) > 0 {
		return db.NewSession(opts...)
	}
	return nil
}
