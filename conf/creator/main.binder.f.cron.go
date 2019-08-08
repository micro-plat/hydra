package creator

type ICronBinder interface {
	imainBinder
}

type CronBinder struct {
	*mainBinder
}

func newCronBinder(params map[string]string, inputs map[string]*Input) *CronBinder {
	return &CronBinder{
		mainBinder: newMainBinder(params, inputs),
	}
}
