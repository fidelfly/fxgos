package db

import (
	"database/sql"

	"github.com/go-xorm/xorm"
)

type DBSession struct {
	orig      *xorm.Session
	autoClose bool
}

func NewSession(params ...bool) *DBSession {
	if Engine == nil {
		panic("database engine is not initialized")
	}
	autoClose := false
	if len(params) > 0 {
		autoClose = params[0]
	}
	return &DBSession{Engine.NewSession(), autoClose}
}

func (dbs *DBSession) GetXorm() *xorm.Session {
	return dbs.orig
}

func (dbs *DBSession) NoAutoTime() {
	dbs.orig.NoAutoTime()
}

func (dbs *DBSession) Insert(data interface{}, opts ...QueryOption) (int64, error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	dbs.orig = attachOption(dbs.orig, opts...)
	return dbs.orig.Insert(data)
}

func (dbs *DBSession) Update(data interface{}, opts ...QueryOption) (int64, error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	dbs.orig = attachOption(dbs.orig, opts...)
	return dbs.orig.Update(data)
}

func (dbs *DBSession) Get(data interface{}, opts ...QueryOption) (bool, error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	dbs.orig = attachOption(dbs.orig, opts...)
	return dbs.orig.Get(data)
}

func (dbs *DBSession) Delete(data interface{}, opts ...QueryOption) (int64, error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	dbs.orig = attachOption(dbs.orig, opts...)
	return dbs.orig.Delete(data)
}

func (dbs *DBSession) Find(data interface{}, opts ...QueryOption) error {
	if dbs.autoClose {
		defer dbs.Close()
	}
	dbs.orig = attachOption(dbs.orig, opts...)
	return dbs.orig.Find(data)
}

func (dbs *DBSession) Exist(data interface{}, opts ...QueryOption) (bool, error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	dbs.orig = attachOption(dbs.orig, opts...)
	return dbs.orig.Exist(data)
}

func (dbs *DBSession) Count(data interface{}, opts ...QueryOption) (int64, error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	dbs.orig = attachOption(dbs.orig, opts...)
	return dbs.orig.Count(data)
}

func (dbs *DBSession) Close() {
	dbs.orig.Close()
}

func (dbs *DBSession) BeginTransaction() error {
	return dbs.orig.Begin()
}

func (dbs *DBSession) EndTransaction(commit bool) error {
	if commit {
		return dbs.orig.Commit()
	}
	return dbs.orig.Rollback()
}

func (dbs *DBSession) Commit() error {
	return dbs.orig.Commit()
}

func (dbs *DBSession) Rollback() error {
	return dbs.orig.Rollback()
}

func (dbs *DBSession) Exec(sqlOrArgs ...interface{}) (sql.Result, error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	return dbs.orig.Exec(sqlOrArgs...)
}
