package db

import (
	"github.com/fidelfly/fxgo/errorx"
)

type UpdateInfo struct {
	Id   int64
	Cols []string
}

type ListInfo struct {
	Results   int64
	Page      int64
	SortField string
	SortOrder string
	Cond      string
}

func GetPagingOption(req ListInfo) []QueryOption {
	opts := make([]QueryOption, 0)
	results := int(req.Results)
	page := int(req.Page)
	sortField := req.SortField
	sortOrder := req.SortOrder

	if results > 0 {
		if page == 0 {
			page = 1
		}
		opts = append(opts, Limit(results, (page-1)*results))
	}
	if len(sortField) > 0 && sortOrder != "false" {
		if sortOrder == "descend" {
			opts = append(opts, Desc(sortField))
		} else {
			opts = append(opts, Asc(sortField))
		}
	}
	return opts
}

var (
	ErrNotExist = errorx.NewError("err.db.record_not_exist", "record not exist")
)
