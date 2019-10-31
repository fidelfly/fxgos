package db

import (
	"database/sql"

	"github.com/go-xorm/xorm"
)

type Session struct {
	orig          *xorm.Session
	autoClose     bool
	inTransaction bool
	callbacks     []TxCallback
}

type TxCallback func(commit bool)

func CommitCallback(f func()) TxCallback {
	return func(commit bool) {
		if commit {
			f()
		}
	}
}

func RollbackCallback(f func()) TxCallback {
	return func(commit bool) {
		if !commit {
			f()
		}
	}
}

func PairCallback(commitCall, rollbackCall func()) TxCallback {
	return func(commit bool) {
		if commit {
			if commitCall != nil {
				commitCall()
			}
		} else {
			if rollbackCall != nil {
				rollbackCall()
			}
		}
	}
}

func (dbs *Session) AddTxCallback(calls ...TxCallback) {
	if dbs.inTransaction {
		dbs.callbacks = append(dbs.callbacks, calls...)
	} else {
		for _, callback := range calls {
			callback(true)
		}
	}
}

func NewSession(opts ...SessionOption) *Session {
	if Engine == nil {
		panic("database engine is not initialized")
	}
	s := &Session{Engine.NewSession(), false, false, nil}
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (dbs *Session) GetXorm() *xorm.Session {
	return dbs.orig
}

func (dbs *Session) NoAutoTime() {
	dbs.orig.NoAutoTime()
}

func (dbs *Session) Insert(data interface{}, opts ...QueryOption) (affected int64, err error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	attachOption(dbs, opts...)
	affected, err = dbs.orig.Insert(data)

	if !dbs.inTransaction {
		dbs.callback(err == nil)
	}
	return
}

func (dbs *Session) Update(data interface{}, opts ...QueryOption) (affected int64, err error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	attachOption(dbs, opts...)
	affected, err = dbs.orig.Update(data)
	if !dbs.inTransaction {
		dbs.callback(err == nil)
	}
	return
}

func (dbs *Session) Get(data interface{}, opts ...QueryOption) (affected bool, err error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	attachOption(dbs, opts...)
	affected, err = dbs.orig.Get(data)
	if !dbs.inTransaction {
		dbs.callback(err == nil)
	}
	return
}

func (dbs *Session) Delete(data interface{}, opts ...QueryOption) (affected int64, err error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	attachOption(dbs, opts...)
	affected, err = dbs.orig.Delete(data)
	if !dbs.inTransaction {
		dbs.callback(err == nil)
	}
	return
}

func (dbs *Session) Find(data interface{}, opts ...QueryOption) (err error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	attachOption(dbs, opts...)
	err = dbs.orig.Find(data)
	if !dbs.inTransaction {
		dbs.callback(err == nil)
	}
	return
}

func (dbs *Session) Exist(data interface{}, opts ...QueryOption) (exist bool, err error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	attachOption(dbs, opts...)
	exist, err = dbs.orig.Exist(data)
	if !dbs.inTransaction {
		dbs.callback(err == nil)
	}
	return
}

func (dbs *Session) Count(data interface{}, opts ...QueryOption) (count int64, err error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	attachOption(dbs, opts...)
	count, err = dbs.orig.Count(data)
	if !dbs.inTransaction {
		dbs.callback(err == nil)
	}
	return
}

func (dbs *Session) Close() {
	if dbs.inTransaction {
		dbs.Rollback()
	}
	dbs.orig.Close()
}

func (dbs *Session) Begin() error {
	if dbs.inTransaction {
		return nil
	}
	dbs.inTransaction = true
	return dbs.orig.Begin()
}

func (dbs *Session) Commit() error {
	if dbs.inTransaction {
		dbs.inTransaction = false
		if err := dbs.orig.Commit(); err != nil {
			return err
		}
		dbs.callback(true)
	}
	return nil
}

func (dbs *Session) Rollback() error {
	if dbs.inTransaction {
		dbs.inTransaction = false
		if err := dbs.orig.Rollback(); err != nil {
			return err
		}
		dbs.callback(false)
	}
	return nil
}

func (dbs *Session) Exec(sqlOrArgs ...interface{}) (result sql.Result, err error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	result, err = dbs.orig.Exec(sqlOrArgs...)
	if !dbs.inTransaction {
		dbs.callback(err == nil)
	}
	return
}

func (dbs *Session) callback(commit bool) {
	if len(dbs.callbacks) > 0 {
		for i := len(dbs.callbacks) - 1; i >= 0; i-- {
			dbs.callbacks[i](commit)
		}
		dbs.callbacks = nil
	}
}
