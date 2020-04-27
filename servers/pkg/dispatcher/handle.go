package dispatcher

type HandlerFunc func(*Context)

func (h HandlerFunc) Handle(ctx *Context) {
	h(ctx)
}

type HandlersChain []HandlerFunc

// Last returns the last handler in the chain. ie. the last handler is the main own.
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}
