package rpc

import "testing"
import "github.com/micro-plat/lib4go/ut"

func TestFactoryResolvePath(t *testing.T) {
	def_domain := "hydra"
	def_server := "sys.api"
	_, svs, domain, server, err := ResolvePath("order.request", def_domain, def_server)
	ut.Expect(t, err, nil)
	ut.Expect(t, svs, "/order/request")
	ut.Expect(t, domain, def_domain)
	ut.Expect(t, server, def_server)

	_, svs, domain, server, err = ResolvePath("/order/request", def_domain, def_server)
	ut.Expect(t, err, nil)
	ut.Expect(t, svs, "/order/request")
	ut.Expect(t, domain, def_domain)
	ut.Expect(t, server, def_server)

	_, svs, domain, server, err = ResolvePath("/order/request@", def_domain, def_server)
	ut.Expect(t, err, nil)
	ut.Expect(t, svs, "/order/request")
	ut.Expect(t, domain, def_domain)
	ut.Expect(t, server, def_server)

	_, svs, domain, server, err = ResolvePath("@", def_domain, def_server)
	ut.Refute(t, err, nil)

	svs, domain, server, err = ResolvePath("@merchant_cron", def_domain, def_server)
	ut.Refute(t, err, nil)

	_, svs, domain, server, err = ResolvePath("/order/request@merchant", def_domain, def_server)
	ut.Expect(t, err, nil)
	ut.Expect(t, svs, "/order/request")
	ut.Expect(t, domain, def_domain)
	ut.Expect(t, server, "merchant")

	_, svs, domain, server, err = ResolvePath("order.request@merchant.", def_domain, def_server)
	ut.Expect(t, err, nil)
	ut.Expect(t, svs, "/order/request")
	ut.Expect(t, domain, def_domain)
	ut.Expect(t, server, "merchant")

	_, svs, domain, server, err = ResolvePath("order.request@merchant.sys", def_domain, def_server)
	ut.Expect(t, err, nil)
	ut.Expect(t, svs, "/order/request")
	ut.Expect(t, domain, "sys")
	ut.Expect(t, server, "merchant")

	_, svs, domain, server, err = ResolvePath("order.request@.sys", def_domain, def_server)
	ut.Expect(t, err, nil)
	ut.Expect(t, svs, "/order/request")
	ut.Expect(t, domain, "sys")
	ut.Expect(t, server, def_server)

	_, svs, domain, server, err = ResolvePath("order/request@merchant.", def_domain, def_server)
	ut.Expect(t, err, nil)
	ut.Expect(t, svs, "/order/request")
	ut.Expect(t, domain, def_domain)
	ut.Expect(t, server, "merchant")

	_, svs, domain, server, err = ResolvePath("order/request@merchant.sys", def_domain, def_server)
	ut.Expect(t, err, nil)
	ut.Expect(t, svs, "/order/request")
	ut.Expect(t, domain, "sys")
	ut.Expect(t, server, "merchant")

	_, svs, domain, server, err = ResolvePath("order/request@.sys", def_domain, def_server)
	ut.Expect(t, err, nil)
	ut.Expect(t, svs, "/order/request")
	ut.Expect(t, domain, "sys")
	ut.Expect(t, server, def_server)

	_, svs, domain, server, err = ResolvePath("/order/request/create@merchant_rpc", def_domain, def_server)
	ut.Expect(t, err, nil)
	ut.Expect(t, svs, "/order/request/create")
	ut.Expect(t, domain, def_domain)
	ut.Expect(t, server, "merchant_rpc")

}
