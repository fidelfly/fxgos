package user

import (
	"context"
	"errors"

	"github.com/fidelfly/gox/pkg/strx"

	"github.com/fidelfly/fxgos/cmd/service/user/res"
	"github.com/fidelfly/fxgos/cmd/utilities/auth"
	"github.com/fidelfly/fxgos/cmd/utilities/mdbo"
	"github.com/fidelfly/fxgos/cmd/utilities/pub"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/gostool/db"
	"github.com/fidelfly/gostool/dbo"
)

func New() *res.User {
	return &res.User{}
}

func Create(ctx context.Context, input interface{}) (user *res.User, err error) {
	if t, ok := input.(*res.User); ok {
		user = t
	} else {
		user = new(res.User)
	}

	err = dbo.Create(ctx,
		dbo.ApplyBeanOption(user,
			dbo.Assignment(input),
			mdbo.CreateUser(ctx),
			dbo.FuncBeanOption(func(target interface{}) {
				if user, ok := target.(*res.User); ok {
					user.Password = auth.EncodePassword(user.Code, user.Password)
				}
			})),
		mdbo.PubResourceEvent(ResourceType, pub.ResourceCreate),
	)
	return
}

type UpdateInput struct {
	db.UpdateInfo
	Data *res.User
}

func Update(ctx context.Context, info dbo.UpdateInfo) error {
	if info.Data == nil {
		return syserr.ErrInvalidParam
	}

	var user *res.User
	if t, ok := info.Data.(*res.User); ok {
		user = t
	} else {
		user = new(res.User)
	}

	opts := dbo.ApplytUpdateOption(user, info, mdbo.UpdateUser(ctx))

	id := user.Id
	pwdChange := len(user.Password) > 0
	roleChange := false
	statusChange := false

	if len(info.Cols) > 0 {
		pwdChange = strx.IndexOfSlice(info.Cols, "password") >= 0
		roleChange = strx.IndexOfSlice(info.Cols, "role_id") >= 0
		statusChange = strx.IndexOfSlice(info.Cols, "status") >= 0
	}

	resUser, err := Read(ctx, id)
	if err != nil {
		return err
	}

	if pwdChange {
		user.Password = auth.EncodePassword(resUser.Code, user.Password)
	}
	if statusChange {
		if resUser.Status != StatusInvalid && user.Status == StatusDeactivated {
			return errors.New("user status is not invalid")
		}
		if resUser.Status != StatusDeactivated && user.Status == StatusValid {
			return errors.New("user status is not deactived")
		}
	}
	if _, err := dbo.Update(ctx, user, opts,
		mdbo.PubResourceEvent(ResourceType, pub.ResourceUpdate),
		dbo.SessionAfter(func(ctx context.Context, bean interface{}) {
			if roleChange {
				pub.Publish(nil, pub.TopicRoleUpdate) //todo add real role update
			}
		}),
	); err != nil {
		return err
	}
	return nil
}

func Read(ctx context.Context, id int64) (*res.User, error) {
	if id <= 0 {
		return nil, syserr.ErrInvalidParam
	}
	user := &res.User{Id: id}
	if find, err := dbo.Read(ctx, user); err != nil {
		return nil, syserr.DatabaseErr(err)
	} else if !find {
		return nil, syserr.ErrNotFound
	}
	return user, nil
}

//export
func ReadByCode(ctx context.Context, code string) (*res.User, error) {
	if len(code) <= 0 {
		return nil, syserr.ErrInvalidParam
	}
	user := &res.User{Code: code}
	if find, err := dbo.Read(ctx, user); err != nil {
		return nil, syserr.DatabaseErr(err)
	} else if !find {
		return nil, syserr.ErrNotFound
	}
	return user, nil
}

func ReadByEmail(ctx context.Context, email string) (*res.User, error) {
	if len(email) <= 0 {
		return nil, syserr.ErrInvalidParam
	}
	user := &res.User{Email: email}
	if find, err := dbo.Read(ctx, user); err != nil {
		return nil, err
	} else if find {
		return user, nil
	}
	return nil, syserr.ErrNotFound
}

func Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return syserr.ErrInvalidParam
	}
	resUser := &res.User{Id: id}
	if find, err := db.Read(resUser); err != nil {
		return syserr.DatabaseErr(err)
	} else if !find {
		return syserr.ErrNotFound
	}
	if resUser.Status != StatusDeleted {
		resUser.Status = StatusDeleted
		if _, err := dbo.Update(ctx, resUser,
			[]db.QueryOption{db.ID(id), db.Cols("status")},
			mdbo.PubResourceEvent(ResourceType, pub.ResourceDelete),
		); err != nil {
			return syserr.DatabaseErr(err)
		}
	}

	return nil
}

type ValidateInput struct {
	Id       int64
	Code     string
	Email    string
	Password string
}

func Validate(ctx context.Context, input ValidateInput) (*res.User, error) {
	user := &res.User{
		Id:    input.Id,
		Code:  input.Code,
		Email: input.Email,
	}
	if ok, _ := db.Read(user, db.Where("status = ?", StatusValid)); ok {
		encodePwd := auth.EncodePassword(user.Code, input.Password)
		if encodePwd == user.Password {
			return user, nil
		}
	}
	return nil, errors.New("invalid user or password")
}

func List(ctx context.Context, input *dbo.ListInfo, conds ...string) (results []*res.User, count int64, err error) {
	results = make([]*res.User, 0)

	count, err = dbo.List(ctx, results, input, db.Condition(conds...))

	if err != nil {
		return nil, 0, err
	}

	return
}

const ResourceType = "resource.user"
