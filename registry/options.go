package registry

import (
	"github.com/micro-plat/lib4go/logger"
)

type Options struct {
	Domain   string
	Addrs    []string
	Auth     *AuthCreds
	Logger   logger.ILogging
	Metadata map[string]string
	PoolSize int
}

type AuthCreds struct {
	Username string
	Password string
}

type Option func(*Options)

// Addrs is the registry addresses to use
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

func WithAuthCreds(username, password string) Option {
	return func(o *Options) {
		o.Auth = &AuthCreds{
			Username: username,
			Password: password,
		}
	}
}

func WithMetadata(i map[string]string) Option {
	return func(o *Options) {
		o.Metadata = i
	}
}

func WithLogger(log logger.ILogging) Option {
	return func(o *Options) {
		o.Logger = log
	}
}

func Metadata(key, val string) Option {
	return func(o *Options) {
		if o.Metadata == nil {
			o.Metadata = map[string]string{}
		}
		o.Metadata[key] = val
	}
}

func WithDomain(domain string) Option {
	return func(o *Options) {
		o.Domain = domain
	}
}
