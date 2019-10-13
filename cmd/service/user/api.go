package user

import (
	"context"
	"errors"

	"github.com/fidelfly/gox/pkg/strh"

	"github.com/fidelfly/fxgos/cmd/service/user/res"
	"github.com/fidelfly/fxgos/cmd/utilities/auth"
	"github.com/fidelfly/fxgos/cmd/utilities/pub"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/gostool/db"
)

func New() *res.User {
	return &res.User{}
}

func Create(ctx context.Context, input *res.User) (id int64, err error) {
	input.Password = auth.EncodePassword(input.Code, input.Password)
	_, err = db.Create(input)
	if err == nil {
		pub.Publish(pub.ResourceEvent{
			Type:   ResourceType,
			Action: pub.ResourceCreate,
			Id:     input.Id,
		}, pub.TopicResource)
	}
	return
}

type UpdateInput struct {
	db.UpdateInfo
	Data *res.User
}

func Update(ctx context.Context, input UpdateInput) (int64, error) {
	if input.Data == nil {
		return 0, errors.New("data is empty")
	}

	id := input.Data.Id
	pwdChange := len(input.Data.Password) > 0
	roleChange := false
	statusChange := false
	opts := make([]db.QueryOption, 0)
	if input.Id > 0 {
		opts = append(opts, db.ID(input.Id))
		id = input.Id
	}
	if len(input.Cols) > 0 {
		opts = append(opts, db.Cols(input.Cols...))
		pwdChange = strh.IndexOfSlice(input.Cols, "password") >= 0
		roleChange = strh.IndexOfSlice(input.Cols, "role_id") >= 0
		statusChange = strh.IndexOfSlice(input.Cols, "status") >= 0
	}

	resUser := &res.User{Id: id}
	if find, err := db.Read(resUser); err != nil {
		return 0, err
	} else if !find {
		return 0, nil
	}
	if pwdChange {
		input.Data.Password = auth.EncodePassword(resUser.Code, input.Data.Password)
	}
	if statusChange {
		if resUser.Status != StatusInvalid && input.Data.Status == StatusDeactivated {
			return 0, errors.New("user status is not invalid")
		}
		if resUser.Status != StatusDeactivated && input.Data.Status == StatusValid {
			return 0, errors.New("user status is not deactived")
		}
	}
	if rows, err := db.Update(input.Data, opts...); err != nil {
		return 0, err
	} else if rows > 0 {
		pub.Publish(pub.ResourceEvent{
			Type:   ResourceType,
			Action: pub.ResourceUpdate,
			Id:     input.Data.Id,
		}, pub.TopicResource)

		if roleChange {
			pub.Publish(nil, pub.TopicRoleUpdate)
			//_ = mss.RolePub.Publish(ctx, &iam.RoleEvent{UserId: resUser.Id, RoleId: resUser.RoleId})
		}
	}
	return input.Id, nil
}

func Read(ctx context.Context, id int64) (*res.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid value of id")
	}
	user := &res.User{Id: id}
	if find, err := db.Read(user); err != nil {
		return nil, err
	} else if find {
		return user, nil
	}
	return nil, nil
}

//export
func ReadByCode(ctx context.Context, code string) (*res.User, error) {
	if len(code) <= 0 {
		return nil, syserr.ErrInvalidParam
	}
	user := &res.User{Code: code}
	if find, err := db.Read(user); err != nil {
		return nil, err
	} else if find {
		return user, nil
	}
	return nil, syserr.ErrNotFound
}

func ReadByEmail(ctx context.Context, email string) (*res.User, error) {
	if len(email) <= 0 {
		return nil, syserr.ErrInvalidParam
	}
	user := &res.User{Email: email}
	if find, err := db.Read(user); err != nil {
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
		if _, err := db.Update(resUser, db.ID(id), db.Cols("status")); err != nil {
			return syserr.DatabaseErr(err)
		}
	}
	pub.Publish(pub.ResourceEvent{
		Type:   ResourceType,
		Action: pub.ResourceDelete,
		Id:     resUser.Id,
	}, pub.TopicResource)

	return nil
	//mskit.RemoveUserCache(resUser.Id)
	//_ = mss.RolePub.Publish(ctx, &iam.RoleEvent{UserId: resUser.Id})
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

func List(ctx context.Context, input db.ListInfo) (results []*res.User, count int64, err error) {
	results = make([]*res.User, 0)
	opts := make([]db.QueryOption, 0)
	if len(input.Cond) > 0 {
		opts = append(opts, db.Where(input.Cond))
	}
	queOpts := append(append(db.GetPagingOption(input), db.Desc("create_time")), opts...)
	if err = db.Find(&results, queOpts...); err != nil {
		return nil, 0, err
	}

	count = int64(len(results))
	if !(count < input.Results && input.Page == 1) {
		count, _ = db.Count(new(res.User), opts...)
	}
	return
}

const ResourceType = "resource.user"
