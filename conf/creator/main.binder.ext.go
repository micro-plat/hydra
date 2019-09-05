package creator

type ExtBinder struct {
	*microBinder
}

func NewExtBinder(params map[string]string, inputs map[string]*Input) *ExtBinder {
	return &ExtBinder{
		microBinder: newMicroBinder(params, inputs),
	}
}
