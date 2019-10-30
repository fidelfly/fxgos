package role

import (
	"context"
	"fmt"

	"github.com/fidelfly/gox/errorx"

	"github.com/fidelfly/fxgos/cmd/service/role/res"
	res2 "github.com/fidelfly/fxgos/cmd/service/user/res"
	"github.com/fidelfly/fxgos/cmd/utilities/mdbo"
	"github.com/fidelfly/fxgos/cmd/utilities/pub"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/gostool/db"
	"github.com/fidelfly/gostool/dbo"
)

type Form struct {
	Code        string  `json:"code"`
	Roles       []int64 `json:"roles"`
	Description string  `json:"description"`
}

func Create(ctx context.Context, input interface{}) (*res.Role, error) {
	var role *res.Role
	if t, ok := input.(*res.Role); ok {
		role = t
	} else {
		role = new(res.Role)
	}

	err := dbo.Create(ctx,
		dbo.ApplyBeanOption(role, dbo.Assignment(input), mdbo.CreateUser(ctx)),
		mdbo.PubResourceEvent(ResourceType, pub.ResourceCreate),
	)
	return role, err
}

func Update(ctx context.Context, info dbo.UpdateInfo) error {
	if info.Data == nil {
		return syserr.ErrInvalidParam
	}

	var role *res.Role
	if t, ok := info.Data.(*res.Role); ok {
		role = t
	} else {
		role = new(res.Role)
	}

	opts := dbo.ApplytUpdateOption(role, info, mdbo.UpdateUser(ctx))

	if rows, err := dbo.Update(ctx, role, opts,
		mdbo.PubResourceEvent(ResourceType, pub.ResourceUpdate)); err != nil {
		return syserr.DatabaseErr(err)
	} else if rows == 0 {
		return syserr.ErrNotFound
	}
	return nil
}

func Read(ctx context.Context, id int64) (*res.Role, error) {
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

func ReadByCode(ctx context.Context, code string) (*res.Role, error) {
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

func Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return syserr.ErrInvalidParam
	}
	resRole := &res.Role{Id: id}
	ctx, dbs := dbo.WithDBSession(ctx, db.AutoClose(false))
	defer dbs.Close()
	if count, err := dbo.Count(ctx, res.Role{}, db.Where("roles like ?", fmt.Sprintf("%%%d%%", id))); err != nil {
		return syserr.DatabaseErr(err)
	} else if count > 0 {
		return errorx.NewError(syserr.CodeOfDatabaseErr, "role_used_by_role")
	}

	if count, err := dbo.Count(ctx, res2.User{}, db.Where("roles like ?", fmt.Sprintf("%%%d%%", id))); err != nil {
		return syserr.DatabaseErr(err)
	} else if count > 0 {
		return errorx.NewError(syserr.CodeOfDatabaseErr, "role_used_by_user")
	}

	if count, err := dbo.Delete(ctx, resRole, nil,
		mdbo.PubResourceEvent(ResourceType, pub.ResourceDelete)); err != nil {
		return syserr.DatabaseErr(err)
	} else if count == 0 {
		return syserr.ErrNotFound
	}

	return nil
}

func List(ctx context.Context, input *dbo.ListInfo, conds ...string) ([]*res.Role, int64, error) {
	resRoles := make([]*res.Role, 0)

	count, err := dbo.List(ctx, resRoles, input, db.Condition(conds...))

	if err != nil {
		return nil, 0, err
	}

	return resRoles, count, nil
}

const ResourceType = "resource.role"
