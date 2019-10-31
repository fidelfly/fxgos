package dbo

import (
	"context"

	"github.com/fidelfly/gostool/db"
)

type SessionHook interface {
	Option(ctx context.Context, dbs *db.Session) //todo should remove???
	Before(ctx context.Context, bean interface{})
	After(ctx context.Context, bean interface{})
}

type smOption struct {
	opts []db.StatementOption
}

func (so *smOption) Option(ctx context.Context, dbs *db.Session) {
	if len(so.opts) > 0 {
		for _, opt := range so.opts {
			opt(dbs)
		}
	}
}
func (so *smOption) Before(ctx context.Context, bean interface{}) {
}
func (so *smOption) After(ctx context.Context, bean interface{}) {
}
func WithStatementOption(opts ...db.StatementOption) SessionHook {
	return &smOption{opts: opts}
}

type sessionBefore func(ctx context.Context, bean interface{})

func (sb sessionBefore) Option(ctx context.Context, dbs *db.Session) {
}
func (sb sessionBefore) Before(ctx context.Context, bean interface{}) {
	sb(ctx, bean)
}
func (sb sessionBefore) After(ctx context.Context, bean interface{}) {
}

func SessionBefore(f func(ctx context.Context, bean interface{})) SessionHook {
	return sessionBefore(f)
}

type sessionAfter func(ctx context.Context, bean interface{})

func (sa sessionAfter) Option(ctx context.Context, dbs *db.Session) {
}
func (sa sessionAfter) Before(ctx context.Context, bean interface{}) {
}
func (sa sessionAfter) After(ctx context.Context, bean interface{}) {
	sa(ctx, bean)
}

func SessionAfter(f func(ctx context.Context, bean interface{})) SessionHook {
	return sessionAfter(f)
}

func applyHooks(ctx context.Context, dbs *db.Session, hooks ...SessionHook) {
	if len(hooks) > 0 {
		for _, hook := range hooks {
			hook.Option(ctx, dbs)
		}

		beforeAction := func(bean interface{}) {
			for _, hook := range hooks {
				hook.Before(ctx, bean)
			}
		}
		afterAction := func(bean interface{}) {
			for _, hook := range hooks {
				hook.After(ctx, bean)
			}
		}
		dbs.GetXorm().Before(beforeAction)
		dbs.GetXorm().After(afterAction)
	}
}
