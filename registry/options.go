package registry

import (
	"crypto/tls"
	"time"

	"github.com/micro-plat/lib4go/logger"
)

type Options struct {
	Domain    string
	Addrs     []string
	Timeout   time.Duration
	TLSConfig *tls.Config
	Auth      *AuthCreds
	Logger    logger.ILogging
	Metadata  map[string]string
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

func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}
func TLSConfig(cfg *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = cfg
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

func Domain(domain string) Option {
	return func(o *Options) {
		o.Domain = domain
	}
}

