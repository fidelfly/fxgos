package db

type SessionOption func(*Session)

func AutoClose(autoClose bool) SessionOption {
	return func(session *Session) {
		session.autoClose = autoClose
	}
}

type StatementOption func(session *Session)

func StatementOptionChain(options []StatementOption) StatementOption {
	return func(session *Session) {
		for _, opt := range options {
			if opt != nil {
				opt(session)
			}
		}
	}
}

func Cols(cols ...string) StatementOption {
	return func(session *Session) {
		session.GetXorm().Cols(cols...)
	}
}

func Table(name string) StatementOption {
	return func(session *Session) {
		session.GetXorm().Table(name)
	}
}

func ID(id interface{}) StatementOption {
	return func(session *Session) {
		session.GetXorm().ID(id)
	}
}

func Where(query interface{}, args ...interface{}) StatementOption {
	return func(session *Session) {
		session.GetXorm().Where(query, args...)
	}
}

func Condition(conds ...string) StatementOption {
	return func(session *Session) {
		for _, cond := range conds {
			session.GetXorm().And(cond)
		}
	}
}

func Limit(limit int, start ...int) StatementOption {
	return func(session *Session) {
		session.GetXorm().Limit(limit, start...)
	}
}

func Asc(colNames ...string) StatementOption {
	return func(session *Session) {
		session.GetXorm().Asc(colNames...)
	}
}

func Desc(colNames ...string) StatementOption {
	return func(session *Session) {
		session.GetXorm().Desc(colNames...)
	}
}

func NoAutoTime() StatementOption {
	return func(session *Session) {
		session.GetXorm().NoAutoTime()
	}
}

func AfterClosure(closure func(interface{})) StatementOption {
	return func(session *Session) {
		session.GetXorm().After(closure)
	}
}

func BeforeClosure(closure func(interface{})) StatementOption {
	return func(session *Session) {
		session.GetXorm().Before(closure)
	}
}

func WithTxCallback(callbacks ...TxCallback) StatementOption {
	return func(session *Session) {
		session.callbacks = append(session.callbacks, callbacks...)
	}
}

func attachOption(session *Session, opts ...StatementOption) {
	if len(opts) > 0 {
		for _, opt := range opts {
			if opt != nil {
				opt(session)
			}
		}
	}
}
