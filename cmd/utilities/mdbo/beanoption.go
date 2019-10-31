package mdbo

import (
	"context"
	"time"

	"github.com/fidelfly/gox/pkg/reflectx"
	"github.com/fidelfly/gox/pkg/strx"

	"github.com/fidelfly/fxgos/cmd/utilities/mctx"
	"github.com/fidelfly/gostool/dbo"

	"github.com/fidelfly/gox/pkg/datex"
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

func ToDayStart(fields ...string) dbo.BeanOption {
	return dbo.FuncBeanOption(func(target interface{}) {
		for _, field := range fields {
			if v := reflectx.GetField(target, field); v != nil {
				if t, ok := v.(time.Time); ok {
					reflectx.SetField(target, reflectx.FV{
						Field: field,
						Value: datex.DateStart(t),
					})
				}
			}
		}
	})
}

func ToDayEnd(fields ...string) dbo.BeanOption {
	return dbo.FuncBeanOption(func(target interface{}) {
		for _, field := range fields {
			if v := reflectx.GetField(target, field); v != nil {
				if t, ok := v.(time.Time); ok {
					reflectx.SetField(target, reflectx.FV{
						Field: field,
						Value: datex.DateEnd(t),
					})
				}
			}
		}
	})
}
