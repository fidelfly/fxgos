package mctx

import (
	"context"

	"github.com/fidelfly/gox/pkg/metax"
)

func GetUserId(ctx context.Context) int64 {
	md := GetCallOption(ctx)
	if v, ok := md.Get(metaCallUser); ok {
		return v.(int64)
	}
	return 0
}

type callOptionKey struct {
}

const (
	metaCallUser     = "meta.callOption.user"
	metaCallAsSA     = "meta.callOption.asSA"
	metaCallIgnoreSg = "meta.callOption.ignoreSg"
)

func WithCallOption(ctx context.Context, options ...metax.MetaOption) context.Context {
	callOption := metax.Wrap(GetCallOption(ctx), metax.ApplyOption(nil, options...))
	return context.WithValue(ctx, callOptionKey{}, callOption)
}

func GetCallOption(ctx context.Context) metax.MetaData {
	if v := ctx.Value(callOptionKey{}); v != nil {
		if md, ok := v.(metax.MetaData); ok {
			return md
		}
	}

	return metax.EmptyMD
}

func CallAsUser(userId int64) metax.MetaOption {
	return func(md metax.MetaData) metax.MetaData {
		_ = md.Set(metaCallUser, userId)
		return md
	}
}

func CallAsSA(md metax.MetaData) metax.MetaData {
	_ = md.Set(metaCallAsSA, true)
	return md
}

func CallIgnoreSg(md metax.MetaData) metax.MetaData {
	_ = md.Set(metaCallIgnoreSg, true)
	return md
}

func IsSuperAdmin(ctx context.Context) bool {
	md := GetCallOption(ctx)
	if v, ok := md.Get(metaCallAsSA); ok {
		return v.(bool)
	}
	return false
}

func IsIgnoreSg(ctx context.Context) bool {
	md := GetCallOption(ctx)
	if v, ok := md.Get(metaCallIgnoreSg); ok {
		return v.(bool)
	}
	return false
}
