package wss

type Frame struct {
	Version  string                 `json:"version"`
	Type     string                 `json:"type"`
	ID       string                 `json:"id,omitempty"`
	Group    string                 `json:"group,omitempty"`
	Client   string                 `json:"client,omitempty"`
	Method   string                 `json:"method,omitempty"`
	Path     string                 `json:"path,omitempty"`
	Query    string                 `json:"query,omitempty"`
	Header   map[string]string      `json:"header,omitempty"`
	Status   int                    `json:"status,omitempty"`
	Error    string                 `json:"error,omitempty"`
	Body     []byte                 `json:"body,omitempty"`
	Services []ServiceMeta          `json:"services,omitempty"`
	Extra    map[string]interface{} `json:"extra,omitempty"`
}

type ServiceMeta struct {
	Path    string   `json:"path"`
	Methods []string `json:"methods,omitempty"`
}

const (
	frameHello       = "hello"
	frameHelloAck    = "hello_ack"
	frameRegister    = "register"
	frameRegistered  = "registered"
	frameRequest     = "request"
	frameResponse    = "response"
	frameStreamStart = "stream_start"
	frameStreamChunk = "stream_chunk"
	frameStreamEnd   = "stream_end"
	framePing        = "ping"
	framePong        = "pong"
	frameError       = "error"
)
