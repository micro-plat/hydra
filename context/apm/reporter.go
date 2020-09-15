package apm

//Reporter is a data transit specification
type Reporter interface {
	Boot(service string, serviceInstance string)
	Send(spans []Span)
	Close()
	GetRealReporter() interface{}
}
