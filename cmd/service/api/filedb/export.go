package filedb

import (
	"context"
	"fmt"

	"github.com/fidelfly/fxgos/cmd/service"
	"github.com/fidelfly/fxgos/cmd/service/res"
)

func getServer() Service {
	if v, ok := service.GetService(ServiceName); ok {
		if server, ok := v.(Service); ok {
			return server
		}
	}
	panic(fmt.Sprintf("Service(%s) is not registered", ServiceName))
}

func Save(ctx context.Context, name string, data []byte) (int64, error) {
	return getServer().Save(ctx, name, data)
}
func Read(ctx context.Context, id int64) (*res.File, error) {
	return getServer().Read(ctx, id)
}
