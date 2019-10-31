package db

type SessionOption func(*Session)

func AutoClose(autoClose bool) SessionOption {
	return func(session *Session) {
		session.autoClose = autoClose
	}
}

type QueryOption func(session *Session)

func Cols(cols ...string) QueryOption {
	return func(session *Session) {
		session.GetXorm().Cols(cols...)
	}
}

func Table(name string) QueryOption {
	return func(session *Session) {
		session.GetXorm().Table(name)
	}
}

func ID(id interface{}) QueryOption {
	return func(session *Session) {
		session.GetXorm().ID(id)
	}
}

func Where(query interface{}, args ...interface{}) QueryOption {
	return func(session *Session) {
		session.GetXorm().Where(query, args...)
	}
}

func Condition(conds ...string) QueryOption {
	return func(session *Session) {
		for _, cond := range conds {
			session.GetXorm().And(cond)
		}
	}
}

func Limit(limit int, start ...int) QueryOption {
	return func(session *Session) {
		session.GetXorm().Limit(limit, start...)
	}
}

func Asc(colNames ...string) QueryOption {
	return func(session *Session) {
		session.GetXorm().Asc(colNames...)
	}
}

func Desc(colNames ...string) QueryOption {
	return func(session *Session) {
		session.GetXorm().Desc(colNames...)
	}
}

func NoAutoTime() QueryOption {
	return func(session *Session) {
		session.GetXorm().NoAutoTime()
	}
}

func AfterClosure(closure func(interface{})) QueryOption {
	return func(session *Session) {
		session.GetXorm().After(closure)
	}
}

func BeforeClosure(closure func(interface{})) QueryOption {
	return func(session *Session) {
		session.GetXorm().Before(closure)
	}
}

func WithTxCallback(callbacks ...TxCallback) QueryOption {
	return func(session *Session) {
		session.callbacks = append(session.callbacks, callbacks...)
	}
}

func attachOption(session *Session, opts ...QueryOption) {
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(session)
		}
	}
}
