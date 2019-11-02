package dbo

import (
	"context"

	"github.com/fidelfly/gostool/db"
)

type dbSessionKey struct{}

type CtxSession struct {
	*db.Session
	controlled bool
	tranOwner  bool
}

func (cs *CtxSession) Begin() error {
	if !cs.Session.InTransaction() {
		if err := cs.Session.Begin(); err != nil {
			return err
		}
		cs.tranOwner = true
	}

	return nil
}

func (cs *CtxSession) Commit() error {
	if cs.Session.InTransaction() {
		if cs.tranOwner {
			if err := cs.Session.Commit(); err != nil {
				return err
			}
			cs.tranOwner = false
		}
	}
	return nil
}

func (cs *CtxSession) Rollback() error {
	if cs.Session.InTransaction() {
		if cs.tranOwner {
			if err := cs.Session.Rollback(); err != nil {
				return err
			}
			cs.tranOwner = false
		}
	}
	return nil
}

func (cs *CtxSession) Close() {
	if cs.controlled {
		cs.Session.Close()
	} else if cs.tranOwner {
		cs.Rollback()
	}
}

func WithDBSession(ctx context.Context, opts ...db.SessionOption) (context.Context, *CtxSession) {
	ctxDbs := CurrentDBSession(ctx)
	if ctxDbs != nil {
		return ctx, ctxDbs
	}
	dbs := db.NewSession(opts...)
	return context.WithValue(ctx, dbSessionKey{}, dbs), &CtxSession{dbs, true, false}
}

func NewCtxSession(ctx context.Context, opts ...db.SessionOption) (context.Context, *CtxSession) {
	dbs := db.NewSession(opts...)
	return context.WithValue(ctx, dbSessionKey{}, dbs), &CtxSession{dbs, true, false}
}

func CurrentDBSession(ctx context.Context, opts ...db.SessionOption) *CtxSession {
	if v := ctx.Value(dbSessionKey{}); v != nil {
		if dbs, ok := v.(*db.Session); ok {
			return &CtxSession{dbs, false, false}
		}
	}
	if len(opts) > 0 {
		return &CtxSession{db.NewSession(opts...), true, false}
	}
	return nil
}

func DefaultSession(session *db.Session) {
	//do nothing
	return
}
