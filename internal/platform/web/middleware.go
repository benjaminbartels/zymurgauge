package web

// Middleware is a func type that wraps around handlers. These funcs or executed before or after the handler.
type Middleware func(Handler) Handler

// wraps the given handler with each given middleware functions in order.
func wrap(mw []Middleware, handler Handler) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
