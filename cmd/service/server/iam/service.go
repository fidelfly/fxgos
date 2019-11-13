package iam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fidelfly/gox/logx"
	"github.com/fidelfly/gox/pkg/mathx"
	"github.com/tidwall/buntdb"

	"github.com/fidelfly/fxgos/cmd/service/api/iam"
	"github.com/fidelfly/fxgos/cmd/service/res"
	model2 "github.com/fidelfly/fxgos/cmd/service/server/iam/model"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/gostool/db"
	"github.com/fidelfly/gostool/dbo"
)

/*//export
func NewAdapter(resType string, data []byte) *model.Adapter {
	a := &model.Adapter{
		ResourceType: resType,
		InitData:     data,
	}
	return a
}*/

type server struct {
}

func (s server) ValidateAccess(userId int64, premise iam.AccessPremise) (bool, error) {
	if len(premise) > 0 {
		for _, ai := range premise {
			if len(ai.Actions) > 0 {
				for _, act := range ai.Actions {
					ok, err := Validate(iam.Query{
						ResourceType: ai.Type,
						UserId:       userId,
						ResourceCode: ai.Code,
						Action:       act,
					})
					if err != nil {
						return false, err
					}

					if !ok {
						return false, nil
					}
				}
			}
		}
	}
	return true, nil
}

func Validate(query iam.Query) (bool, error) {
	var sub, obj, act string

	if query.RoleId > 0 {
		sub = iam.EncodeRoleSubject(query.RoleId)
	} else if query.UserId > 0 {
		sub = iam.EncodeUserSubject(query.UserId)
	} else {
		return false, errors.New("invalid subject")
	}

	keys := make([]interface{}, len(query.ResourceKey))
	if len(query.ResourceKey) > 0 {
		for i, key := range query.ResourceKey {
			keys[i] = key
		}
	}
	obj = iam.EncodeObject(query.ResourceType, query.ResourceCode, keys...)

	act = query.Action

	return model2.Validate(query.ResourceType, sub, obj, act), nil
}

func (s server) ListResource(ctx context.Context, resourceType string) []*iam.Resource {
	resList := make([]*iam.Resource, 0)
	_ = resDB.GetDB().Update(func(tx *buntdb.Tx) error {
		_ = tx.AscendRange("type",
			fmt.Sprintf(`{"type":"%s", "index": %d, "init_index": %d }`, resourceType, 0, 0),
			fmt.Sprintf(`{"type":"%s", "index": %d, "init_index": %d }`, resourceType, mathx.MaxInt, mathx.MaxInt),
			func(key, value string) bool {
				iamRes := &iam.Resource{}
				if err := json.Unmarshal([]byte(value), iamRes); err == nil {
					resList = append(resList, iamRes)
				}
				return true
			})
		return nil
	})

	return resList
}

func (s server) QueryByUser(ctx context.Context, userId int64, resourceType string) []*iam.AccessItem {
	return s.queryBySubject(ctx, iam.EncodeUserSubject(userId), resourceType)
}

func (s server) QueryByRole(ctx context.Context, roleId int64, resourceType string) []*iam.AccessItem {
	return s.queryBySubject(ctx, iam.EncodeRoleSubject(roleId), resourceType)
}

func (s server) queryBySubject(ctx context.Context, sub string, resourceType string) []*iam.AccessItem {
	resources := s.ListResource(ctx, resourceType)
	if len(resources) == 0 {
		return nil
	}
	var obj, act string
	items := make([]*iam.AccessItem, 0)
	for _, iamRes := range resources {
		obj = iam.EncodeObject(iamRes.Type, iamRes.Code)
		if len(iamRes.Actions) > 0 {
			policyAction := make([]string, 0)
			for _, act = range iamRes.Actions {
				if model2.Validate(iamRes.Type, sub, obj, act) {
					policyAction = append(policyAction, act)
				}
			}
			if len(policyAction) > 0 {
				items = append(items, &iam.AccessItem{
					Type:    iamRes.Type,
					Code:    iamRes.Code,
					Actions: policyAction,
				})
			}
		}
	}
	return items
}

func (s server) listResourceAcl(ctx context.Context, sub string, resourceType string) []*iam.ResourceACL {
	resources := s.ListResource(ctx, resourceType)
	if len(resources) == 0 {
		return nil
	}

	var obj, act string
	items := make([]*iam.ResourceACL, 0)
	for _, iamRes := range resources {
		item := &iam.ResourceACL{
			Resource: *iamRes,
		}
		policyAction := make([]string, 0)
		if len(sub) > 0 {
			obj = iam.EncodeObject(iamRes.Type, iamRes.Code)
			if len(iamRes.Actions) > 0 {
				for _, act = range iamRes.Actions {
					if model2.Validate(iamRes.Type, sub, obj, act) {
						policyAction = append(policyAction, act)
					}
				}
			}
		}
		item.ACL = policyAction
		items = append(items, item)
	}
	return items
}

func (s server) ListResourceAclByRole(ctx context.Context, roleId int64, resourctType string) []*iam.ResourceACL {
	return s.listResourceAcl(ctx, iam.EncodeRoleSubject(roleId), resourctType)
}

func (s server) ListResourceAclByUser(ctx context.Context, userId int64, resourctType string) []*iam.ResourceACL {
	return s.listResourceAcl(ctx, iam.EncodeUserSubject(userId), resourctType)
}

func (s server) UpdatePolicyByRole(ctx context.Context, roleId int64, inheritRoles []int64, acl []*iam.ResourceACL) (err error) {
	policies := make([]*res.Policy, 0)
	for _, item := range acl {
		if len(item.ACL) > 0 {
			newPolicy := &res.Policy{
				RoleId:       roleId,
				ResourceType: item.Type,
				Sub:          iam.EncodeRoleSubject(roleId),
				Obj:          iam.EncodeObject(item.Type, item.Code),
				Act:          item.ACL,
			}
			policies = append(policies, newPolicy)
			/*		for _, action := range item.ACL {
					newPolicy := &res.Policy{
						RoleId:       roleId,
						ResourceType: item.Type,
						Sub:          iamx.EncodeRoleSubject(roleId),
						Obj:          iamx.EncodeObject(item.Type, item.Code),
						Act:          action,
					}
					policies = append(policies, newPolicy)
				}*/
		}
	}

	dbSession := dbo.CurrentDBSession(ctx, dbo.DefaultSession)
	defer dbSession.Close()
	if err := dbSession.Begin(); err != nil {
		return syserr.DatabaseErr(err)
	}
	defer func() {
		if err != nil {
			logx.CaptureError(dbSession.Rollback())
		}
	}()
	if _, err := dbSession.Exec("delete from policy where role_id = ? ", roleId); err != nil {
		return syserr.DatabaseErr(err)
	}

	if len(policies) > 0 {
		if _, err := dbSession.Insert(&policies); err != nil {
			return syserr.DatabaseErr(err)
		}
	}

	if err := dbSession.Commit(); err != nil {
		return syserr.DatabaseErr(err)
	}

	dbSession.AddTxCallback(db.CommitCallback(func() {
		model2.UpdateRolePolicy(roleId, inheritRoles...)
	}))

	return nil
}

func (s server) UpdatePolicyByUser(ctx context.Context, userId int64, roles []int64, superAdmin bool) error {
	model2.UpdateUserRole(userId, roles, superAdmin)
	return nil
}

func (s server) DeleteRolePolicy(ctx context.Context, roleId int64) error {
	if rows, err := db.Delete(&res.Policy{RoleId: roleId}); err != nil {
		return syserr.DatabaseErr(err)
	} else if rows == 0 {
		return syserr.ErrNotFound
	}
	model2.DeleteRolePolicy(roleId)
	return nil
}

func (s server) DeleteUserPolicy(ctx context.Context, userId int64) error {
	model2.DeleteUserPolicy(userId)
	return nil
}
