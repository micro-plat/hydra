package apm

const (
	TagURL             = "url"
	TagStatusCode      = "status_code"
	TagHTTPMethod      = "http.method"
	TagDBType          = "db.type"
	TagDBInstance      = "db.instance"
	TagDBStatement     = "db.statement"
	TagDBBindVariables = "db.bind_vars"
	TagMQQueue         = "mq.queue"
	TagMQBroker        = "mq.broker"
	TagMQTopic         = "mq.topic"
)

const (
	SpanLayer_Unknown      = 0
	SpanLayer_Database     = 1
	SpanLayer_RPCFramework = 2
	SpanLayer_Http         = 3
	SpanLayer_MQ           = 4
	SpanLayer_Cache        = 5
)

var (
	Header string = "sw8"
)

//SpanOption SpanOption
type SpanOption func(s Span)

//WithPeer x
func WithPeer(peer string) SpanOption {
	return func(s Span) {
		s.SetPeer(peer)
	}
}

// WithSpanLayer x
func WithSpanLayer(spanLayer int32) SpanOption {
	return func(s Span) {
		s.SetSpanLayer(spanLayer)
	}
}

//WithTag x
func WithTag(k, v string) SpanOption {
	return func(s Span) {
		s.Tag(k, v)
	}
}
