package role

import (
	"context"
	"fmt"

	"github.com/fidelfly/fxgos/cmd/service"
	"github.com/fidelfly/fxgos/cmd/service/res"
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

func Create(ctx context.Context, input interface{}) (*res.Role, error) {
	return getServer().Create(ctx, input)
}
func Update(ctx context.Context, info dbo.UpdateInfo) error {
	return getServer().Update(ctx, info)
}
func Read(ctx context.Context, id int64) (*res.Role, error) {
	return getServer().Read(ctx, id)
}
func ReadByCode(ctx context.Context, code string) (*res.Role, error) {
	return getServer().ReadByCode(ctx, code)
}
func Delete(ctx context.Context, id int64) error {
	return getServer().Delete(ctx, id)
}
func List(ctx context.Context, input *dbo.ListInfo, conds ...string) ([]*res.Role, int64, error) {
	return getServer().List(ctx, input, conds...)
}
