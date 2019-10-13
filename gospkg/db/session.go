package db

import (
	"database/sql"

	"github.com/go-xorm/xorm"
)

type Session struct {
	orig      *xorm.Session
	autoClose bool
}

func NewSession(params ...bool) *Session {
	if Engine == nil {
		panic("database engine is not initialized")
	}
	autoClose := false
	if len(params) > 0 {
		autoClose = params[0]
	}
	return &Session{Engine.NewSession(), autoClose}
}

func (dbs *Session) GetXorm() *xorm.Session {
	return dbs.orig
}

func (dbs *Session) NoAutoTime() {
	dbs.orig.NoAutoTime()
}

func (dbs *Session) Insert(data interface{}, opts ...QueryOption) (int64, error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	dbs.orig = attachOption(dbs.orig, opts...)
	return dbs.orig.Insert(data)
}

func (dbs *Session) Update(data interface{}, opts ...QueryOption) (int64, error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	dbs.orig = attachOption(dbs.orig, opts...)
	return dbs.orig.Update(data)
}

func (dbs *Session) Get(data interface{}, opts ...QueryOption) (bool, error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	dbs.orig = attachOption(dbs.orig, opts...)
	return dbs.orig.Get(data)
}

func (dbs *Session) Delete(data interface{}, opts ...QueryOption) (int64, error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	dbs.orig = attachOption(dbs.orig, opts...)
	return dbs.orig.Delete(data)
}

func (dbs *Session) Find(data interface{}, opts ...QueryOption) error {
	if dbs.autoClose {
		defer dbs.Close()
	}
	dbs.orig = attachOption(dbs.orig, opts...)
	return dbs.orig.Find(data)
}

func (dbs *Session) Exist(data interface{}, opts ...QueryOption) (bool, error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	dbs.orig = attachOption(dbs.orig, opts...)
	return dbs.orig.Exist(data)
}

func (dbs *Session) Count(data interface{}, opts ...QueryOption) (int64, error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	dbs.orig = attachOption(dbs.orig, opts...)
	return dbs.orig.Count(data)
}

func (dbs *Session) Close() {
	dbs.orig.Close()
}

func (dbs *Session) BeginTransaction() error {
	return dbs.orig.Begin()
}

func (dbs *Session) EndTransaction(commit bool) error {
	if commit {
		return dbs.orig.Commit()
	}
	return dbs.orig.Rollback()
}

func (dbs *Session) Commit() error {
	return dbs.orig.Commit()
}

func (dbs *Session) Rollback() error {
	return dbs.orig.Rollback()
}

func (dbs *Session) Exec(sqlOrArgs ...interface{}) (sql.Result, error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	return dbs.orig.Exec(sqlOrArgs...)
}
