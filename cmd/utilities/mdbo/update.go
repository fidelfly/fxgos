package mdbo

import (
	"context"

	"github.com/fidelfly/gox/pkg/reflectx"
	"github.com/fidelfly/gox/pkg/strx"

	"github.com/fidelfly/fxgos/cmd/utilities/mctx"
	"github.com/fidelfly/gostool/db"
	"github.com/fidelfly/gostool/dbo"
)

func UpdateUser(ctx context.Context, fields ...string) dbo.UpdateOption {
	userId := mctx.GetUserId(ctx)
	if userId == 0 {
		return nil
	}
	if strx.IndexOfSlice(fields, "UpdateUser") < 0 {
		fields = append(fields, "UpdateUser")
	}
	return dbo.FuncUpdateOption(func(target interface{}) []db.QueryOption {
		pairs := make([]reflectx.FV, len(fields))
		for i, f := range fields {
			pairs[i] = reflectx.FV{
				Field: f,
				Value: userId,
			}
		}
		fields := reflectx.SetField(target, pairs...)
		if len(fields) > 0 {
			cols := make([]string, len(fields))
			for i, v := range fields {
				cols[i] = strx.UnderscoreString(v)
			}
			return []db.QueryOption{db.Cols(cols...)}
		}
		return nil
	})
}
