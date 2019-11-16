package iam

import (
	"context"

	"github.com/fidelfly/fxgos/cmd/service"
)

const (
	ServiceName = "service.iam"
)

func RegisterServer(server Service, dependencies ...string) {
	service.Register(ServiceName, server, dependencies...)
}

type Service interface {
	ListResource(ctx context.Context, resourceType string) []*Resource

	ValidateAccess(userId int64, premise AccessPremise) (bool, error)
	QueryByUser(ctx context.Context, userId int64, resourceType string) []*AccessItem
	QueryByRole(ctx context.Context, roleId int64, resourceType string) []*AccessItem

	ListResourceAclByRole(ctx context.Context, roleId int64, resourctType string) []*ResourceACL
	ListResourceAclByUser(ctx context.Context, userId int64, resourctType string) []*ResourceACL

	UpdatePolicyByRole(ctx context.Context, roleId int64, inheritRoles []int64, acl []*ResourceACL) (err error)
	UpdatePolicyByUser(ctx context.Context, userId int64, roles []int64, superAdmin bool) error
	DeleteRolePolicy(ctx context.Context, roleId int64) error
	DeleteUserPolicy(ctx context.Context, userId int64) error
}

type Query struct {
	ResourceType string
	UserId       int64
	RoleId       int64
	ResourceCode string
	ResourceKey  []string
	Action       string
}

type Resource struct {
	Type    string   `json:"type"`
	Code    string   `json:"code"`
	Actions []string `json:"actions"`
}

type ResourceACL struct {
	Resource
	ACL []string `json:"acl"`
}
