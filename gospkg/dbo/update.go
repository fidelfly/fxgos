package dbo

import (
	"context"

	"github.com/fidelfly/gox/pkg/reflectx"
	"github.com/fidelfly/gox/pkg/strx"

	"github.com/fidelfly/gostool/db"
)

type UpdateOption interface {
	Apply(interface{}) []db.StatementOption
}

type FuncUpdateOption func(target interface{}) []db.StatementOption

func (fuo FuncUpdateOption) Apply(target interface{}) []db.StatementOption {
	return fuo(target)
}

func ApplyUpdateOption(target interface{}, option ...UpdateOption) []db.StatementOption {
	queryOption := make([]db.StatementOption, 0)
	for _, opt := range option {
		if !reflectx.IsValueNil(opt) {
			if qopt := opt.Apply(target); len(qopt) > 0 {
				queryOption = append(queryOption, qopt...)
			}
		}
	}
	return queryOption
}

func UpdateField(s interface{}, fields ...string) UpdateOption {
	return FuncUpdateOption(func(t interface{}) []db.StatementOption {
		updateFields := reflectx.CopyFields(t, s, fields...)
		if len(updateFields) > 0 {
			fields := make([]string, len(updateFields))
			for i, v := range updateFields {
				fields[i] = strx.UnderscoreString(v)
			}
			return []db.StatementOption{db.Cols(fields...)}
		}
		return nil
	})
}

type UpdateInfo struct {
	Id   int64
	Cols []string
	Data interface{}
}

func (info UpdateInfo) Apply(target interface{}) []db.StatementOption {
	options := make([]db.StatementOption, 0)
	if info.Id > 0 {
		reflectx.SetField(target, reflectx.FV{Field: "Id", Value: info.Id})
		options = append(options, db.ID(info.Id))
	} else if info.Data != nil {
		if v := reflectx.GetField(info.Data, "Id"); v != nil {
			if id, ok := v.(int64); ok {
				reflectx.SetField(target, reflectx.FV{Field: "Id", Value: id})
				options = append(options, db.ID(info.Id))
			}
		}
	}
	if info.Data != nil && info.Data != target {
		reflectx.CopyFields(target, info.Data, info.Cols...)
	}
	if len(info.Cols) > 0 {
		options = append(options, db.Cols(info.Cols...))
	}
	return options
}

func Update(ctx context.Context, target interface{}, option ...db.StatementOption) (int64, error) {
	dbs := CurrentDBSession(ctx, db.AutoClose(true))
	if effectRows, err := dbs.Update(target, option...); err != nil {
		return 0, err
	} else {
		return effectRows, nil
	}
}
