package dbo

import (
	"context"

	"github.com/fidelfly/gox/pkg/ctxx"
	"github.com/fidelfly/gox/pkg/reflectx"

	"github.com/fidelfly/fxgos/cmd/utilities/pub"
)

func getId(target interface{}) int64 {
	if v := reflectx.GetField(target, "Id"); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

func getCode(target interface{}) string {
	if v := reflectx.GetField(target, "Code"); v != nil {
		if code, ok := v.(string); ok {
			return code
		}
	}
	return ""
}

func GetCURDMeta(ctx context.Context) *CURDMeta {
	return &CURDMeta{md: ctxx.GetMetadata(ctx)}
}

func PubResourceEvent(resourceType string, action int) SessionHook {
	return SessionAfter(func(ctx context.Context, bean interface{}) {
		pub.Publish(pub.ResourceEvent{
			Type:   resourceType,
			Action: action,
			Id:     getId(bean),
			Code:   getCode(bean),
		}, pub.TopicResource)
	})
}

func Assignment(s interface{}) BeanOption {
	return func(t interface{}) {
		reflectx.CopyFields(t, s)
	}
}
