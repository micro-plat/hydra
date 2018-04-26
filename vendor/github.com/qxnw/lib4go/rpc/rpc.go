package rpc

import "time"

type IRPCResponse interface {
	GetService() string
	Wait(timeout time.Duration) (int, string, map[string]string, error)
	GetResult() chan IRPCResult
}
type IRPCResult interface {
	GetService() string
	GetStatus() int
	GetResult() string
	GetParams() map[string]string
	GetErr() error
}
