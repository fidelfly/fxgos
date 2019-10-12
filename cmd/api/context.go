package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/fidelfly/fxgo/gosrvx"
	"github.com/fidelfly/fxgo/httprxr"
	"github.com/fidelfly/fxgo/pkg/ctxh"

	"github.com/fidelfly/fxgos/cmd/service/user"
)

const (
	AccessConfigKey       = "access.config"
	ProgressSubscriberKey = "context.progress.subscriber"
)

func GetUserInfo(r *http.Request) *user.CacheInfo {
	userKey := GetUserKey(r)
	if len(userKey) > 0 {
		id, err := strconv.ParseInt(userKey, 10, 64)
		if err != nil {
			return nil
		}
		return user.GetCache(id)
	}
	return nil
}

func GetUserKey(r *http.Request) string {
	userKey := httprxr.ContextGet(r, gosrvx.ContextUserKey)
	if userKey == nil {
		return ""
	}
	if key, ok := userKey.(string); ok {
		return key
	}
	return ""
}

func GetRequestId(r *http.Request) string {
	reqId := httprxr.ContextGet(r, gosrvx.ContextRequestId)
	if reqId == nil {
		return ""
	}
	if key, ok := reqId.(string); ok {
		return key
	}
	return ""
}

func ProgressSubscriber(subscribers ...httprxr.ProgressSubscriber) ctxh.ContextValueHolder {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ProgressSubscriberKey, subscribers)
	}
}

func GetProgressSubscriber(ctx context.Context) []httprxr.ProgressSubscriber {
	ctxObj := ctx.Value(ProgressSubscriberKey)
	if ctxObj == nil {
		return nil
	}

	if subscribers, ok := ctxObj.([]httprxr.ProgressSubscriber); ok {
		return subscribers
	}
	return nil
}
