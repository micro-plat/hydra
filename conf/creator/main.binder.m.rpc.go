package creator

type IRPCBinder interface {
	imicroBinder
}

type RpcBinder struct {
	*microBinder
}

func NewRpcBinder(params map[string]string, inputs map[string]*Input) *RpcBinder {
	return &RpcBinder{
		microBinder: newMicroBinder(params, inputs),
	}
}
