package iam

import (
	"context"
	"fmt"

	"github.com/fidelfly/fxgos/cmd/service"
)

type AccessPremise []AccessItem

type AccessItem struct {
	Type    string
	Code    string
	Keys    []string
	Actions []string
}

//export
func NewAccessItem(resType string, code string, actions ...string) AccessItem {
	return AccessItem{
		Type:    resType,
		Code:    code,
		Actions: actions,
	}
}

func getServer() Service {
	if v, ok := service.GetService(ServiceName); ok {
		if server, ok := v.(Service); ok {
			return server
		}
	}
	panic(fmt.Sprintf("Service(%s) is not registered", ServiceName))
}

func ListResource(ctx context.Context, resourceType string) []*Resource {
	return getServer().ListResource(ctx, resourceType)
}

func ValidateAccess(userId int64, premise AccessPremise) (bool, error) {
	return getServer().ValidateAccess(userId, premise)
}
func QueryByUser(ctx context.Context, userId int64, resourceType string) []*AccessItem {
	return getServer().QueryByUser(ctx, userId, resourceType)
}
func QueryByRole(ctx context.Context, roleId int64, resourceType string) []*AccessItem {
	return getServer().QueryByRole(ctx, roleId, resourceType)
}

func ListResourceAclByRole(ctx context.Context, roleId int64, resourctType string) []*ResourceACL {
	return getServer().ListResourceAclByRole(ctx, roleId, resourctType)
}
func ListResourceAclByUser(ctx context.Context, userId int64, resourctType string) []*ResourceACL {
	return getServer().ListResourceAclByUser(ctx, userId, resourctType)
}

func UpdatePolicyByRole(ctx context.Context, roleId int64, inheritRoles []int64, acl []*ResourceACL) (err error) {
	return getServer().UpdatePolicyByRole(ctx, roleId, inheritRoles, acl)
}
func UpdatePolicyByUser(ctx context.Context, userId int64, roles []int64, superAdmin bool) error {
	return getServer().UpdatePolicyByUser(ctx, userId, roles, superAdmin)
}
func DeleteRolePolicy(ctx context.Context, roleId int64) error {
	return getServer().DeleteRolePolicy(ctx, roleId)
}
func DeleteUserPolicy(ctx context.Context, userId int64) error {
	return getServer().DeleteUserPolicy(ctx, userId)
}
