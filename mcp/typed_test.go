package mcp

import (
	"testing"

	"github.com/micro-plat/hydra/context"
)

type tReq struct{ Name string }
type tResp struct{ Msg string }

func fnTypedValue(ctx context.IContext, req tReq) (tResp, error)  { return tResp{Msg: req.Name}, nil }
func fnTypedPtr(ctx context.IContext, req *tReq) (*tResp, error)  { return &tResp{Msg: req.Name}, nil }
func fnOld(ctx context.IContext) interface{}                      { return nil }

// TestWrapTypedRecognize 验证 typed Handle 签名识别：仅 func(IContext, req)(resp, error) 被识别，
// 旧签名与非函数/参数个数不符的不识别（旧签名由 swapFunc 前两条分支处理，不走 typed hook）。
func TestWrapTypedRecognize(t *testing.T) {
	cases := []struct {
		name string
		fn   interface{}
		want bool
	}{
		{"typed-value", fnTypedValue, true},
		{"typed-ptr", fnTypedPtr, true},
		{"old-signature", fnOld, false},
		{"non-func", 123, false},
		{"one-in", func(a int) int { return a }, false},
		{"three-in", func(a, b, c int) (int, error) { return 0, nil }, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, ok := wrapTyped(c.fn)
			if ok != c.want {
				t.Fatalf("wrapTyped(%s) ok=%v want %v", c.name, ok, c.want)
			}
		})
	}
}
