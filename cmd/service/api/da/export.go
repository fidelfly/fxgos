package da

import (
	"context"
	"fmt"
	"strings"

	"github.com/fidelfly/fxgos/cmd/service"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/fxgos/cmd/utilities/mctx"
	"github.com/fidelfly/gostool/db"
	"github.com/fidelfly/gostool/dbo"
)

func getServer() Service {
	if v, ok := service.GetService(ServiceName); ok {
		if server, ok := v.(Service); ok {
			return server
		}
	}
	panic(fmt.Sprintf("Service(%s) is not registered", ServiceName))
}

func RegisterData(data SecurityData) error {
	return getServer().RegisterData(data)
}
func ValidateAccess(userId int64, resourceType string, resId int64) bool {
	return getServer().ValidateAccess(userId, resourceType, resId)
}
func Create(ctx context.Context, input interface{}) (*res.SecurityGroup, error) {
	return getServer().Create(ctx, input)
}
func Update(ctx context.Context, info dbo.UpdateInfo) error {
	return getServer().Update(ctx, info)
}
func Read(ctx context.Context, id int64) (*res.SecurityGroup, error) {
	return getServer().Read(ctx, id)
}
func ReadByCode(ctx context.Context, code string) (*res.SecurityGroup, error) {
	return getServer().ReadByCode(ctx, code)
}
func Delete(ctx context.Context, id int64) error {
	return getServer().Delete(ctx, id)
}
func List(ctx context.Context, input *dbo.ListInfo, conds ...string) ([]*res.SecurityGroup, int64, error) {
	return getServer().List(ctx, input, conds...)
}

func SgCondition(ctx context.Context, resourceType string, tableAlias string) string {
	if mctx.IsIgnoreSg(ctx) || mctx.IsSuperAdmin(ctx) {
		return ""
	}
	userId := mctx.GetUserId(ctx)
	return strings.Join([]string{
		"exists (",
		"select 1 from user_sg as sg1, resource_sg as sg2",
		fmt.Sprintf(`where sg2.res_type = "%s" and sg2.res_id = %s.id`, resourceType, tableAlias),
		fmt.Sprintf(`and sg1.user_id=%d`, userId),
		"and sg1.security_group = sg2.security_group",
		")",
	}, " ")
}

func SgCondOption(ctx context.Context, resourceType string, tableAlias string) db.StatementOption {
	cond := SgCondition(ctx, resourceType, tableAlias)
	if len(cond) > 0 {
		return db.Where(cond)
	}
	return nil
}
