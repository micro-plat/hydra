package registry

// Watcher is an interface that returns updates
// about services within the registry.
type Watcher interface {
	// Next is a blocking call
	Next() (*Result, error)
	Stop()
}

type Result struct {
	Action string
	Data   interface{}
}

type WatchOptions struct {
	// Specify a service to watch
	// If blank, the watch is for all services
	Service string

	// Domain to watch
	Domain string
}

type WatchOption func(*WatchOptions)

func WatchService(name string) WatchOption {
	return func(o *WatchOptions) {
		o.Service = name
	}
}

func WatchDomain(d string) WatchOption {
	return func(o *WatchOptions) {
		o.Domain = d
	}
}
