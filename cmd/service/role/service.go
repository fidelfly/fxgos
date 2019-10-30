package role

import (
	"context"

	"github.com/fidelfly/fxgos/cmd/service/role/res"
	"github.com/fidelfly/fxgos/cmd/utilities/mctx"
	"github.com/fidelfly/fxgos/cmd/utilities/pub"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/gostool/db"
)

func _create(ctx context.Context, resRole *res.Role) (int64, error) {
	if resRole == nil {
		return 0, syserr.ErrInvalidParam
	}
	mctx.FillUserInfo(ctx, resRole)
	if _, err := db.Create(resRole); err != nil {
		return 0, syserr.DatabaseErr(err)
	} else {
		pub.Publish(pub.ResourceEvent{
			Type:   ResourceType,
			Action: pub.ResourceCreate,
			Id:     resRole.Id,
		}, pub.TopicResource)
		return resRole.Id, nil
	}
}
