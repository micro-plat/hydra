package creator

type IMQCBinder interface {
	imainBinder
}

type MQCBinder struct {
	*mainBinder
}

func NewMQCBinder(params map[string]string, inputs map[string]*Input) *MQCBinder {
	return &MQCBinder{
		mainBinder: newMainBinder(params, inputs),
	}
}
