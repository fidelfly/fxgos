package role

import (
	"context"
	"fmt"

	"github.com/fidelfly/gox/errorx"

	"github.com/fidelfly/fxgos/cmd/service/api/role"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/fxgos/cmd/utilities/mdbo"
	"github.com/fidelfly/fxgos/cmd/utilities/pub"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/gostool/db"
	"github.com/fidelfly/gostool/dbo"
)

type server struct {
}

func (s server) Create(ctx context.Context, input interface{}) (*res.Role, error) {
	if input == nil {
		return nil, syserr.ErrInvalidParam
	}
	var resRole *res.Role
	if t, ok := input.(*res.Role); ok {
		resRole = t
	} else {
		resRole = new(res.Role)
	}

	err := dbo.Create(ctx,
		dbo.ApplyBeanOption(resRole, dbo.Assignment(input), mdbo.CreateUser(ctx)),
		mdbo.ResourceEventHook(role.ResourceType, pub.ResourceCreate),
	)
	return resRole, err
}

func (s server) Update(ctx context.Context, info dbo.UpdateInfo) error {
	if info.Data == nil {
		return syserr.ErrInvalidParam
	}

	var resRole *res.Role
	if t, ok := info.Data.(*res.Role); ok {
		resRole = t
	} else {
		resRole = new(res.Role)
	}

	opts := dbo.ApplyUpdateOption(resRole, info, mdbo.UpdateUser(ctx))

	if rows, err := dbo.Update(ctx, resRole,
		db.StatementOptionChain(opts),
		mdbo.ResourceEventHook(role.ResourceType, pub.ResourceUpdate)); err != nil {
		return syserr.DatabaseErr(err)
	} else if rows == 0 {
		return syserr.ErrNotFound
	}
	return nil
}

func (s server) Read(ctx context.Context, id int64) (*res.Role, error) {
	if id <= 0 {
		return nil, syserr.ErrInvalidParam
	}
	resRole := &res.Role{Id: id}
	if find, err := dbo.Read(ctx, resRole); err != nil {
		return nil, syserr.DatabaseErr(err)
	} else if !find {
		return nil, syserr.ErrNotFound
	}
	return resRole, nil
}

func (s server) ReadByCode(ctx context.Context, code string) (*res.Role, error) {
	if len(code) == 0 {
		return nil, syserr.ErrInvalidParam
	}
	resRole := &res.Role{Code: code}
	if find, err := dbo.Read(ctx, resRole); err != nil {
		return nil, syserr.DatabaseErr(err)
	} else if !find {
		return nil, syserr.ErrNotFound
	}
	return resRole, nil
}

func (s server) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return syserr.ErrInvalidParam
	}
	resRole := &res.Role{Id: id}
	dbs := dbo.CurrentDBSession(ctx, dbo.DefaultSession)
	defer dbs.Close()
	if count, err := dbs.Count(res.Role{}, db.Where("roles like ?", fmt.Sprintf("%%%d%%", id))); err != nil {
		return syserr.DatabaseErr(err)
	} else if count > 0 {
		return errorx.NewError(syserr.CodeOfDatabaseErr, "role_used_by_role")
	}

	if count, err := dbs.Count(res.User{}, db.Where("roles like ?", fmt.Sprintf("%%%d%%", id))); err != nil {
		return syserr.DatabaseErr(err)
	} else if count > 0 {
		return errorx.NewError(syserr.CodeOfDatabaseErr, "role_used_by_user")
	}

	if count, err := dbs.Delete(resRole,
		mdbo.ResourceEventHook(role.ResourceType, pub.ResourceDelete)); err != nil {
		return syserr.DatabaseErr(err)
	} else if count == 0 {
		return syserr.ErrNotFound
	}

	return nil
}

func (s server) List(ctx context.Context, input *dbo.ListInfo, conds ...string) ([]*res.Role, int64, error) {
	resRoles := make([]*res.Role, 0)

	count, err := dbo.List(ctx, &resRoles, input, db.Condition(conds...))

	if err != nil {
		return nil, 0, err
	}

	return resRoles, count, nil
}
