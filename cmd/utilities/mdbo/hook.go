package mdbo

import (
	"context"

	"github.com/fidelfly/fxgos/cmd/utilities/pub"
	"github.com/fidelfly/gostool/db"
	"github.com/fidelfly/gostool/dbo"
)

func ResourceEventHook(resourceType string, action int) dbo.SessionHook {
	return dbo.SessionAfter(func(ctx context.Context, bean interface{}) {
		pub.Publish(pub.ResourceEvent{
			Type:   resourceType,
			Action: action,
			Id:     getId(bean),
			Code:   getCode(bean),
		}, pub.TopicResource)
	})
}

func ResourceEventOption(resourceType string, action int) db.StatementOption {
	return db.AfterClosure(func(bean interface{}) {
		pub.Publish(pub.ResourceEvent{
			Type:   resourceType,
			Action: action,
			Id:     getId(bean),
			Code:   getCode(bean),
		}, pub.TopicResource)
	})
}
