package da

import (
	"context"

	"github.com/fidelfly/fxgos/cmd/service/api/da"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/fxgos/cmd/utilities/mdbo"
	"github.com/fidelfly/fxgos/cmd/utilities/pub"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/gostool/db"
	"github.com/fidelfly/gostool/dbo"
)

type server struct {
}

func (s server) ValidateAccess(userId int64, resourceType string, resId int64) bool {
	return validateDataAccess(encodeUserKey(userId), resource{resourceType, resId})
}

func (s server) Create(ctx context.Context, input interface{}) (*res.SecurityGroup, error) {
	if input == nil {
		return nil, syserr.ErrInvalidParam
	}
	var resSg *res.SecurityGroup
	if t, ok := input.(*res.SecurityGroup); ok {
		resSg = t
	} else {
		resSg = new(res.SecurityGroup)
	}

	err := dbo.Create(ctx,
		dbo.ApplyBeanOption(resSg, dbo.Assignment(input), mdbo.CreateUser(ctx)),
		mdbo.ResourceEventHook(da.ResourceType, pub.ResourceCreate),
	)
	return resSg, err
}
func (s server) Update(ctx context.Context, info dbo.UpdateInfo) error {
	if info.Data == nil {
		return syserr.ErrInvalidParam
	}

	var resSg *res.SecurityGroup
	if t, ok := info.Data.(*res.SecurityGroup); ok {
		resSg = t
	} else {
		resSg = new(res.SecurityGroup)
	}

	opts := dbo.ApplyUpdateOption(resSg, info, mdbo.UpdateUser(ctx))

	if rows, err := dbo.Update(ctx, resSg,
		db.StatementOptionChain(opts),
		mdbo.ResourceEventHook(da.ResourceType, pub.ResourceUpdate)); err != nil {
		return syserr.DatabaseErr(err)
	} else if rows == 0 {
		return syserr.ErrNotFound
	}
	return nil
}
func (s server) Read(ctx context.Context, id int64) (*res.SecurityGroup, error) {
	if id <= 0 {
		return nil, syserr.ErrInvalidParam
	}
	resSg := &res.SecurityGroup{Id: id}
	if find, err := dbo.Read(ctx, resSg); err != nil {
		return nil, syserr.DatabaseErr(err)
	} else if !find {
		return nil, syserr.ErrNotFound
	}
	return resSg, nil
}
func (s server) ReadByCode(ctx context.Context, code string) (*res.SecurityGroup, error) {
	if len(code) == 0 {
		return nil, syserr.ErrInvalidParam
	}
	resSg := &res.SecurityGroup{Code: code}
	if find, err := dbo.Read(ctx, resSg); err != nil {
		return nil, syserr.DatabaseErr(err)
	} else if !find {
		return nil, syserr.ErrNotFound
	}
	return resSg, nil
}

func (s server) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return syserr.ErrInvalidParam
	}
	resSg := &res.SecurityGroup{Id: id}
	dbs := dbo.CurrentDBSession(ctx, dbo.DefaultSession)
	defer dbs.Close()

	if count, err := dbs.Delete(resSg,
		mdbo.ResourceEventHook(da.ResourceType, pub.ResourceDelete)); err != nil {
		return syserr.DatabaseErr(err)
	} else if count == 0 {
		return syserr.ErrNotFound
	}

	return nil
}
func (s server) List(ctx context.Context, input *dbo.ListInfo, conds ...string) ([]*res.SecurityGroup, int64, error) {
	resSgs := make([]*res.SecurityGroup, 0)

	count, err := dbo.List(ctx, &resSgs, input, db.Condition(conds...))

	if err != nil {
		return nil, 0, err
	}

	return resSgs, count, nil
}
