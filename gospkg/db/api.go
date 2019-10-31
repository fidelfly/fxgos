package db

import (
	"database/sql"

	"github.com/go-xorm/xorm"
)

var Engine *xorm.Engine

//export
func Synchronize(beans ...interface{}) error {
	if Engine == nil {
		panic("database engine is not initialized")
	}
	return Engine.Sync2(beans...)
}

func Create(data interface{}, opts ...StatementOption) (int64, error) {
	dbs := NewSession(AutoClose(true))
	return dbs.Insert(data, opts...)
}

func Update(data interface{}, opts ...StatementOption) (int64, error) {
	dbs := NewSession(AutoClose(true))
	return dbs.Update(data, opts...)
}

func Read(data interface{}, opts ...StatementOption) (bool, error) {
	dbs := NewSession(AutoClose(true))
	return dbs.Get(data, opts...)
}

func Delete(data interface{}, opts ...StatementOption) (int64, error) {
	dbs := NewSession(AutoClose(true))
	return dbs.Delete(data, opts...)
}

func Find(data interface{}, opts ...StatementOption) error {
	dbs := NewSession(AutoClose(true))
	return dbs.Find(data, opts...)
}

func Exist(data interface{}, opts ...StatementOption) (bool, error) {
	dbs := NewSession(AutoClose(true))
	return dbs.Exist(data, opts...)
}

func Count(data interface{}, opts ...StatementOption) (int64, error) {
	dbs := NewSession(AutoClose(true))
	return dbs.Count(data, opts...)
}

func Exec(sqlOrArgs ...interface{}) (sql.Result, error) {
	dbs := NewSession(AutoClose(true))
	return dbs.Exec(sqlOrArgs...)
}
