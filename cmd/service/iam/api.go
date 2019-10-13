package iam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fidelfly/gox/logx"
	"github.com/tidwall/buntdb"

	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/service/iam/iamx"
	"github.com/fidelfly/fxgos/cmd/service/iam/model"
	"github.com/fidelfly/fxgos/cmd/service/iam/res"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
)

//export
func NewAdapter(resType string, data []byte) *model.Adapter {
	a := &model.Adapter{
		ResourceType: resType,
		InitData:     data,
	}
	return a
}

//export
func NewAccessItem(resType string, code string, actions ...string) AccessItem {
	return AccessItem{
		Type:    resType,
		Code:    code,
		Actions: actions,
	}
}

type AccessPremise []AccessItem

type AccessItem struct {
	Type    string
	Code    string
	Keys    []string
	Actions []string
}

func ValidateAccess(userId int64, premise AccessPremise) (bool, error) {
	if len(premise) > 0 {
		for _, ai := range premise {
			if len(ai.Actions) > 0 {
				for _, act := range ai.Actions {
					ok, err := Validate(Query{
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

type Query struct {
	ResourceType string
	UserId       int64
	RoleId       int64
	ResourceCode string
	ResourceKey  []string
	Action       string
}

func Validate(query Query) (bool, error) {
	var sub, obj, act string

	if query.RoleId > 0 {
		sub = iamx.EncodeRoleSubject(query.RoleId)
	} else if query.UserId > 0 {
		sub = iamx.EncodeUserSubject(query.UserId)
	} else {
		return false, errors.New("invalid subject")
	}

	keys := make([]interface{}, len(query.ResourceKey))
	if len(query.ResourceKey) > 0 {
		for i, key := range query.ResourceKey {
			keys[i] = key
		}
	}
	obj = iamx.EncodeObject(query.ResourceType, query.ResourceCode, keys...)

	act = query.Action

	return model.Validate(query.ResourceType, sub, obj, act), nil
}

func ListResource(ctx context.Context, resourceType string) []*Resource {
	resList := make([]*Resource, 0)
	_ = resDB.GetDB().Update(func(tx *buntdb.Tx) error {
		_ = tx.AscendEqual("type", fmt.Sprintf(`{"type":"%s"}`, resourceType), func(key, value string) bool {
			iamRes := &Resource{}
			if err := json.Unmarshal([]byte(value), iamRes); err == nil {
				resList = append(resList, iamRes)
			}
			return true
		})
		return nil
	})

	return resList
}

func QueryByUser(ctx context.Context, userId int64, resourceType string) []*AccessItem {
	return queryBySubject(ctx, iamx.EncodeUserSubject(userId), resourceType)
}

func QueryByRole(ctx context.Context, roleId int64, resourceType string) []*AccessItem {
	return queryBySubject(ctx, iamx.EncodeRoleSubject(roleId), resourceType)
}

func queryBySubject(ctx context.Context, sub string, resourceType string) []*AccessItem {
	resources := ListResource(ctx, resourceType)
	if len(resources) == 0 {
		return nil
	}
	var obj, act string
	items := make([]*AccessItem, 0)
	for _, iamRes := range resources {
		obj = iamx.EncodeObject(iamRes.Type, iamRes.Code)
		if len(iamRes.Actions) > 0 {
			policyAction := make([]string, 0)
			for _, act = range iamRes.Actions {
				if model.Validate(iamRes.Type, sub, obj, act) {
					policyAction = append(policyAction, act)
				}
			}
			if len(policyAction) > 0 {
				items = append(items, &AccessItem{
					Type:    iamRes.Type,
					Code:    iamRes.Code,
					Actions: policyAction,
				})
			}
		}
	}
	return items
}

func listResourceAcl(ctx context.Context, sub string, resourceType string) []*ResourceACL {
	resources := ListResource(ctx, resourceType)
	if len(resources) == 0 {
		return nil
	}

	var obj, act string
	items := make([]*ResourceACL, 0)
	for _, iamRes := range resources {
		item := &ResourceACL{
			Resource: *iamRes,
		}
		policyAction := make([]string, 0)
		if len(sub) > 0 {
			obj = iamx.EncodeObject(iamRes.Type, iamRes.Code)
			if len(iamRes.Actions) > 0 {
				for _, act = range iamRes.Actions {
					if model.Validate(iamRes.Type, sub, obj, act) {
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

func ListResourceAclByRole(ctx context.Context, roleId int64, resourctType string) []*ResourceACL {
	return listResourceAcl(ctx, iamx.EncodeRoleSubject(roleId), resourctType)
}

func ListResourceAclByUser(ctx context.Context, userId int64, resourctType string) []*ResourceACL {
	return listResourceAcl(ctx, iamx.EncodeUserSubject(userId), resourctType)
}

func UpdatePolicyByRole(ctx context.Context, roleId int64, inheritRoles []int64, acl []*ResourceACL) (err error) {
	policies := make([]*res.Policy, 0)
	for _, item := range acl {
		if len(item.ACL) > 0 {
			newPolicy := &res.Policy{
				RoleId:       roleId,
				ResourceType: item.Type,
				Sub:          iamx.EncodeRoleSubject(roleId),
				Obj:          iamx.EncodeObject(item.Type, item.Code),
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

	dbSession := db.NewSession()
	defer dbSession.Close()
	if err := dbSession.BeginTransaction(); err != nil {
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
	model.UpdateRolePolicy(roleId, inheritRoles...)
	return nil
}

func UpdatePolicyByUser(ctx context.Context, userId int64, roles []int64) error {
	model.UpdateUserRole(userId, roles)
	return nil
}

func DeleteRolePolicy(ctx context.Context, roleId int64) error {
	if rows, err := db.Delete(&res.Policy{RoleId: roleId}); err != nil {
		return syserr.DatabaseErr(err)
	} else if rows == 0 {
		return syserr.ErrNotFound
	}
	model.DeleteRolePolicy(roleId)
	return nil
}

func DeleteUserPolicy(ctx context.Context, userId int64) error {
	model.DeleteUserPolicy(userId)
	return nil
}
