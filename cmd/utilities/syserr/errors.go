package syserr

import "github.com/fidelfly/gox/errorx"

var (
	ErrNotFound     = errorx.NewError("err.db.not_exist", "data is not found in db")
	ErrInvalidParam = errorx.NewError("err.param.invalid", "invalid params")
)

const (
	CodeOfDatabaseErr = "err.db.error"
	CodeOfServerErr   = "err.server.error"
	CodeOfNotAllowed  = "err.action.not_allowed"
)

func ServerErr(err error) error {
	return errorx.NewCodeError(err, CodeOfServerErr)
}

func DatabaseErr(err error) error {
	return errorx.NewCodeError(err, CodeOfDatabaseErr)
}

func Error(code string, err error) error {
	return errorx.NewCodeError(err, code)
}

func NewError(code string, message string) error {
	return errorx.NewError(code, message)
}
