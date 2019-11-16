package user

import (
	"context"

	"github.com/fidelfly/fxgos/cmd/service/api/da"
	"github.com/fidelfly/fxgos/cmd/service/api/user"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/gostool/db"
)

var s = &server{}

func (s server) Start() error {
	return da.RegisterData(&serverKit{&s})
}

func Initialize() error {
	if err := db.Synchronize(
		new(res.User),
	); err != nil {
		return err
	}

	user.RegisterServer(s, da.ServiceName)
	return nil
}

type serverKit struct {
	srv *server
}

func (sk serverKit) GetResourceType() string {
	return user.ResourceType
}

func (sk serverKit) GetSecurityGroups(ctx context.Context, id int64) ([]int64, error) {
	if resUser, err := sk.srv.Read(ctx, id); err != nil {
		return nil, err
	} else {
		return resUser.SecurityGroups, nil
	}
}
