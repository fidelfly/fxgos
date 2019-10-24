package db

type SessionOption func(*Session)

func AutoClose(autoClose bool) SessionOption {
	return func(session *Session) {
		session.autoClose = autoClose
	}
}

func WithCallback(callbacks ...SessionCallback) SessionOption {
	return func(session *Session) {
		session.callbacks = append(session.callbacks, callbacks...)
	}
}
