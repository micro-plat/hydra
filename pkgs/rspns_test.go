package pkgs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRPCResponse(t *testing.T) {
	//text := `{"channel_no":"testdownchannel","failed_code":"000","failed_msg":"订单成功","order_id":21013630,"request_no":"202002281414911911","status":"SUCCESS"}`

	text := `"{\"channel_no\":\"testdownchannel\",\"failed_code\":\"000\",\"failed_msg\":\"订单成功\",\"order_id\":21013630,\"request_no\":\"202002281414911911\",\"status\":\"SUCCESS\"}"`
	resp := NewRspnsByHD(200, ``, text)
	assert.Equal(t, nil, resp.err)
	assert.Equal(t, 6, len(resp.GetMap()))
}
