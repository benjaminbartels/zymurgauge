package web

// MiddlewareFunc is a func type that wraps around handlers. These funcs or executed before or after the handler.
type MiddlewareFunc func(Handler) Handler

// wraps the given handler with each given middleware function in orde.
func wrap(handler Handler, mw []MiddlewareFunc) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		if mw[i] != nil {
			handler = mw[i](handler)
		}
	}

	return handler
}
