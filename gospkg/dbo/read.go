package dbo

import (
	"context"
	"errors"
	"reflect"

	"github.com/fidelfly/gostool/db"
)

type QueryOption interface {
	Apply() []db.QueryOption
}

type FuncQueryOption func() []db.QueryOption

func (fqo FuncQueryOption) Apply() []db.QueryOption {
	return fqo()
}

type DirectQuery []db.QueryOption

func (dq DirectQuery) Apply() []db.QueryOption {
	return []db.QueryOption(dq)
}

type ListInfo struct {
	Results   int
	Page      int
	SortField string
	SortOrder string
}

func (info ListInfo) Apply() []db.QueryOption {
	opts := make([]db.QueryOption, 0)
	results := info.Results
	page := info.Page
	sortField := info.SortField
	sortOrder := info.SortOrder

	if results > 0 {
		if page == 0 {
			page = 1
		}
		opts = append(opts, db.Limit(results, (page-1)*results))
	}
	if len(sortField) > 0 && sortOrder != "false" {
		if sortOrder == "descend" {
			opts = append(opts, db.Desc(sortField))
		} else {
			opts = append(opts, db.Asc(sortField))
		}
	}
	return opts
}

func Read(ctx context.Context, target interface{}, option ...db.QueryOption) (bool, error) {
	dbs := CurrentDBSession(ctx, db.AutoClose(true))
	return dbs.Get(target, option...)
}

func ApplyQueryOption(option ...QueryOption) []db.QueryOption {
	options := make([]db.QueryOption, 0)
	for _, opt := range option {
		if opt != nil {
			if qos := opt.Apply(); len(qos) > 0 {
				options = append(options, qos...)
			}
		}
	}
	return options
}

func List(ctx context.Context, target interface{}, info *ListInfo, option ...db.QueryOption) (int64, error) {
	if reflect.TypeOf(target).Kind() != reflect.Slice {
		return 0, errors.New("target is not a slice")
	}
	dbs := CurrentDBSession(ctx, db.AutoClose(true))
	if err := dbs.Find(target, ApplyQueryOption(DirectQuery(option), info)...); err != nil {
		return 0, err
	}
	count := reflect.ValueOf(target).Len()
	if info != nil {
		if !(count < info.Results && info.Page == 1) {
			typ := reflect.TypeOf(target).Elem()
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			var countTarget interface{}
			if typ.Kind() == reflect.Struct {
				countTarget = reflect.New(typ).Interface()
			} else {
				countTarget = make(map[string]interface{})
			}
			return dbs.Count(countTarget, option...)
		}
	}

	return int64(count), nil
}

func Find(ctx context.Context, target interface{}, option ...db.QueryOption) error {
	if reflect.TypeOf(target).Kind() != reflect.Slice {
		return errors.New("target is not a slice")
	}
	dbs := CurrentDBSession(ctx, db.AutoClose(true))
	return dbs.Find(target, option...)
}

func Exist(ctx context.Context, target interface{}, option ...db.QueryOption) (bool, error) {
	dbs := CurrentDBSession(ctx, db.AutoClose(true))
	return dbs.Exist(target, option...)
}

func Count(ctx context.Context, target interface{}, option ...db.QueryOption) (int64, error) {
	dbs := CurrentDBSession(ctx, db.AutoClose(true))
	return dbs.Count(target, option...)
}
