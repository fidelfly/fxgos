package mdbo

import (
	"context"

	"github.com/fidelfly/gox/pkg/reflectx"
	"github.com/fidelfly/gox/pkg/strx"

	"github.com/fidelfly/fxgos/cmd/utilities/mctx"
	"github.com/fidelfly/gostool/dbo"
)

func CreateUser(ctx context.Context, fields ...string) dbo.BeanOption {
	userId := mctx.GetUserId(ctx)
	if userId == 0 {
		return nil
	}
	if strx.IndexOfSlice(fields, "CreateUser") < 0 {
		fields = append(fields, "CreateUser")
	}
	if strx.IndexOfSlice(fields, "UpdateUser") < 0 {
		fields = append(fields, "UpdateUser")
	}
	return dbo.FuncBeanOption(func(target interface{}) {
		pairs := make([]reflectx.FV, len(fields))
		for i, f := range fields {
			pairs[i] = reflectx.FV{
				Field: f,
				Value: userId,
			}
		}
		reflectx.SetField(target, pairs...)
	})
}
