package db

import (
	"github.com/fidelfly/gox/errorx"
)

type UpdateInfo struct {
	Id   int64
	Cols []string
}

var (
	ErrNotExist = errorx.NewError("err.db.record_not_exist", "record not exist")
)
