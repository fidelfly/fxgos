package mdbo

import (
	"context"

	"github.com/fidelfly/fxgos/cmd/utilities/pub"
	"github.com/fidelfly/gostool/dbo"
)

func PubResourceEvent(resourceType string, action int) dbo.SessionHook {
	return dbo.SessionAfter(func(ctx context.Context, bean interface{}) {
		pub.Publish(pub.ResourceEvent{
			Type:   resourceType,
			Action: action,
			Id:     getId(bean),
			Code:   getCode(bean),
		}, pub.TopicResource)
	})
}
