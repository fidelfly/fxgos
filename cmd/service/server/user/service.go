package user

import (
	"context"
	"errors"

	"github.com/fidelfly/gox/pkg/strx"

	"github.com/fidelfly/fxgos/cmd/service/api/user"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/fxgos/cmd/utilities/auth"
	"github.com/fidelfly/fxgos/cmd/utilities/mdbo"
	"github.com/fidelfly/fxgos/cmd/utilities/pub"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/gostool/db"
	"github.com/fidelfly/gostool/dbo"
)

type server struct {
}

func (s server) Create(ctx context.Context, input interface{}) (resUser *res.User, err error) {
	if input == nil {
		return nil, syserr.ErrInvalidParam
	}

	if t, ok := input.(*res.User); ok {
		resUser = t
	} else {
		resUser = new(res.User)
	}

	err = dbo.Create(ctx,
		dbo.ApplyBeanOption(resUser,
			dbo.Assignment(input),
			mdbo.CreateUser(ctx),
			dbo.FuncBeanOption(func(target interface{}) {
				if resUser, ok := target.(*res.User); ok {
					resUser.Password = auth.EncodePassword(resUser.Code, resUser.Password)
				}
			})),
		mdbo.ResourceEventHook(user.ResourceType, pub.ResourceCreate),
	)
	return
}

func (s server) Update(ctx context.Context, info dbo.UpdateInfo) error {
	if info.Data == nil {
		return syserr.ErrInvalidParam
	}
	ctx, dbs := dbo.WithDBSession(ctx, dbo.DefaultSession)
	defer dbs.Close()
	var resUser *res.User
	if t, ok := info.Data.(*res.User); ok {
		resUser = t
	} else {
		resUser = new(res.User)
	}

	opts := dbo.ApplyUpdateOption(resUser, info, mdbo.UpdateUser(ctx))

	id := resUser.Id
	orgUser, err := s.Read(ctx, id)
	if err != nil {
		return err
	}

	pwdChange := len(resUser.Password) > 0
	statusChange := resUser.Status != orgUser.Status

	if len(info.Cols) > 0 {
		pwdChange = strx.IndexOfSlice(info.Cols, "password") >= 0
		statusChange = strx.IndexOfSlice(info.Cols, "status") >= 0
	}

	if pwdChange {
		resUser.Password = auth.EncodePassword(orgUser.Code, resUser.Password)
	}
	if statusChange {
		if orgUser.Status != user.StatusInvalid && resUser.Status == user.StatusDeactivated {
			return errors.New("resUser status is not invalid")
		}
		if orgUser.Status != user.StatusDeactivated && resUser.Status == user.StatusValid {
			return errors.New("resUser status is not deactived")
		}
	}
	if _, err := dbo.Update(ctx, resUser, db.StatementOptionChain(opts),
		mdbo.ResourceEventHook(user.ResourceType, pub.ResourceUpdate),
	); err != nil {
		return err
	}
	return nil
}

func (s server) Read(ctx context.Context, id int64) (*res.User, error) {
	if id <= 0 {
		return nil, syserr.ErrInvalidParam
	}
	resUser := &res.User{Id: id}
	if find, err := dbo.Read(ctx, resUser); err != nil {
		return nil, syserr.DatabaseErr(err)
	} else if !find {
		return nil, syserr.ErrNotFound
	}
	return resUser, nil
}

//export
func (s server) ReadByCode(ctx context.Context, code string) (*res.User, error) {
	if len(code) <= 0 {
		return nil, syserr.ErrInvalidParam
	}
	resUser := &res.User{Code: code}
	if find, err := dbo.Read(ctx, resUser); err != nil {
		return nil, syserr.DatabaseErr(err)
	} else if !find {
		return nil, syserr.ErrNotFound
	}
	return resUser, nil
}

func (s server) ReadByEmail(ctx context.Context, email string) (*res.User, error) {
	if len(email) <= 0 {
		return nil, syserr.ErrInvalidParam
	}
	resUser := &res.User{Email: email}
	if find, err := dbo.Read(ctx, resUser); err != nil {
		return nil, err
	} else if find {
		return resUser, nil
	}
	return nil, syserr.ErrNotFound
}

func (s server) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return syserr.ErrInvalidParam
	}
	resUser := &res.User{Id: id}
	dbs := dbo.CurrentDBSession(ctx, dbo.DefaultSession)
	defer dbs.Close()
	if find, err := dbs.Get(resUser); err != nil {
		return syserr.DatabaseErr(err)
	} else if !find {
		return syserr.ErrNotFound
	}
	if resUser.Status != user.StatusDeleted {
		resUser.Status = user.StatusDeleted
		if _, err := dbs.Update(resUser,
			db.ID(id), db.Cols("status"),
			mdbo.ResourceEventHook(user.ResourceType, pub.ResourceDelete),
		); err != nil {
			return syserr.DatabaseErr(err)
		}
	}

	return nil
}

func (s server) Validate(ctx context.Context, input user.ValidateInput) (*res.User, error) {
	resUser := &res.User{
		Id:    input.Id,
		Code:  input.Code,
		Email: input.Email,
	}
	if ok, _ := dbo.Read(ctx, resUser, db.Where("status = ?", user.StatusValid)); ok {
		encodePwd := auth.EncodePassword(resUser.Code, input.Password)
		if encodePwd == resUser.Password {
			return resUser, nil
		}
	}
	return nil, errors.New("invalid user or password")
}

func (s server) List(ctx context.Context, input *dbo.ListInfo, conds ...string) (results []*res.User, count int64, err error) {
	results = make([]*res.User, 0)

	count, err = dbo.List(ctx, &results, input, db.Condition(conds...))

	if err != nil {
		return nil, 0, err
	}

	return
}
