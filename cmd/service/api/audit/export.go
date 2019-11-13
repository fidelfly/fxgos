package audit

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

//export
func ListTrail(ctx context.Context, input *dbo.ListInfo, conds ...string) ([]*res.Systrail, int64, error) {
	return getServer().ListTrail(ctx, input, conds...)
}
