package wss

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/micro-plat/lib4go/types"
)

func isStreamFrame(frame *Frame) bool {
	if frame == nil || frame.Extra == nil {
		return false
	}
	return types.GetBool(frame.Extra["stream"])
}

func wantsStream(r *http.Request, body []byte) bool {
	if r == nil {
		return false
	}
	if strings.Contains(strings.ToLower(r.Header.Get("Accept")), "text/event-stream") {
		return true
	}
	if strings.EqualFold(r.Header.Get("X-Hydra-WSS-Stream"), "true") {
		return true
	}
	if strings.EqualFold(r.URL.Query().Get("stream"), "true") {
		return true
	}
	if len(body) == 0 {
		return false
	}
	var payload map[string]interface{}
	if json.Unmarshal(body, &payload) != nil {
		return false
	}
	return types.GetBool(payload["stream"])
}
