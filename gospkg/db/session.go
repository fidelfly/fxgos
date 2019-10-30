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
	if dbs.inTransaction {
		dbs.Rollback()
	}
	dbs.orig.Close()
}

func (dbs *Session) Begin() error {
	dbs.inTransaction = true
	return dbs.orig.Begin()
}

/*func (dbs *Session) EndTransaction(commit bool) error {
	if commit {
		return dbs.Commit()
	}
	return dbs.Rollback()
}*/

func (dbs *Session) Commit() error {
	dbs.inTransaction = false
	if err := dbs.orig.Commit(); err != nil {
		return err
	}
	dbs.callback(true)
	return nil
}

func (dbs *Session) Rollback() error {
	dbs.inTransaction = false
	if err := dbs.orig.Rollback(); err != nil {
		return err
	}
	dbs.callback(false)
	return nil
}

func (dbs *Session) Exec(sqlOrArgs ...interface{}) (sql.Result, error) {
	if dbs.autoClose {
		defer dbs.Close()
	}
	return dbs.orig.Exec(sqlOrArgs...)
}

func (dbs *Session) callback(commit bool) {
	if len(dbs.callbacks) > 0 {
		for i := len(dbs.callbacks) - 1; i >= 0; i-- {
			dbs.callbacks[i](commit)
		}
		dbs.callbacks = nil
	}
}
