package mdbo

import (
	"github.com/fidelfly/fxgos/cmd/utilities/pub"
	"github.com/fidelfly/gostool/db"
)

func ResourceEventHook(resourceType string, action int) db.StatementOption {
	return db.AfterClosure(func(bean interface{}) {
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
