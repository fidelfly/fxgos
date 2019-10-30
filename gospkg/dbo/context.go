package dbo

import (
	"context"

	"github.com/fidelfly/gostool/db"
)

type dbSessionKey struct{}

type CtxSession struct {
	*db.Session
	controlled bool
}

func (cs *CtxSession) Begin() error {
	if cs.controlled {
		if err := cs.Session.Begin(); err != nil {
			return err
		}
	}
	return nil
}

func (cs *CtxSession) Commit() error {
	if cs.controlled {
		if err := cs.Session.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func (cs *CtxSession) Rollback() error {
	if cs.controlled {
		if err := cs.Session.Rollback(); err != nil {
			return err
		}
	}
	return nil
}

func (cs *CtxSession) Close() {
	if cs.controlled {
		cs.Session.Close()
	}
}

func WithDBSession(ctx context.Context, opts ...db.SessionOption) (context.Context, *CtxSession) {
	ctxDbs := CurrentDBSession(ctx)
	if ctxDbs != nil {
		return ctx, ctxDbs
	}
	dbs := db.NewSession(opts...)
	return context.WithValue(ctx, dbSessionKey{}, dbs), &CtxSession{dbs, true}
}

func NewCtxSession(ctx context.Context, opts ...db.SessionOption) (context.Context, *CtxSession) {
	dbs := db.NewSession(opts...)
	return context.WithValue(ctx, dbSessionKey{}, dbs), &CtxSession{dbs, true}
}

func CurrentDBSession(ctx context.Context, opts ...db.SessionOption) *CtxSession {
	if v := ctx.Value(dbSessionKey{}); v != nil {
		if dbs, ok := v.(*db.Session); ok {
			return &CtxSession{dbs, false}
		}
	}
	if len(opts) > 0 {
		return &CtxSession{db.NewSession(opts...), true}
	}
	return nil
}

func DefaultSession(session *db.Session) {
	//do nothing
	return
}
