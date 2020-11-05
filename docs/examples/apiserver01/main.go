package main

import (
	"fmt"
	nethttp "net/http"
	"time"

	_ "github.com/mattn/go-oci8"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/components"
 	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/types"

	_ "net/http/pprof"
	_ "runtime/pprof"
)

//服务器各种返回结果
func main() {

	go func() {
		nethttp.ListenAndServe("localhost:6060", nil)
	}()

	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithUsage("apiserver"),
		hydra.WithDebug(),
		hydra.WithPlatName("test"),
		hydra.WithSystemName("apiserver"),
	)

	app.API("/order/request/:tp", request)
	app.API("/order/encoding/:tp", request, router.WithEncoding("gbk"))
	app.Start()
}

func request(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("------request------")
	switch ctx.Request().Param("tp") { //从路由配置中获取参数值 ctx.Request.Param.Get...
	case "1":
		time.Sleep(5 * time.Second)

		url := "http://192.168.5.108:8082/order/request"
		client, err := components.Def.HTTP().GetClient()
		if err != nil {
			return fmt.Errorf("GetClient():%+v", err)
		}
		status, content, err := client.Get(url)
		return types.XMap{
			"status":  status,
			"content": content,
			"err":     err,
		}
	case "2":
		return 100
	case "3":
		return time.Now().String()
	case "4":
		return float32(100.20)
	case "5":
		return true
	case "6":
		type order struct {
			ID string `json:"id"`
		}
		type result struct {
			Name   string   `json:"name"`
			Age    int      `json:"age"`
			Orders []*order `json:"order"`
		}
		return &result{Name: "colin", Age: 8, Orders: []*order{&order{ID: "897776666"}}}
	case "7":
		return `{"name":"colin","age":8}`
	case "8":
		return map[string]string{
			"order": "123456",
		}
	case "9":
		return map[string]interface{}{
			"product": map[string]string{
				"price": "100",
			},
		}
	case "10":
		return `<?xml version='1.0'?><xml><name>colin</name><age>8</age></xml>`
	case "11":
		ctx.Response().ContentType("application/xml; charset=UTF-8")
		type order struct {
			ID string `json:"id" xml:"id"`
		}
		type result struct {
			Name   string   `json:"name" xml:"name"`
			Age    int      `json:"age" xml:"age"`
			Orders []*order `json:"orders" xml:"orders"`
		}
		return &result{Name: "colin", Age: 8, Orders: []*order{&order{ID: "897776666"}}}
	case "12":
		return errs.NewError(201, "无需处理")
	case "13":
		if err := ctx.Request().Check("order_id"); err != nil {
			return err
		}
		return ctx.Request().GetString("order_id")
	case "14":
		ctx.Log().Info(ctx.Request().GetBody())
		return fmt.Sprintf(`<?xml version='1.0'?><xml><name>%s</name><age>8</age></xml>`, ctx.Request().GetString("name"))
	case "15":
		ctx.Log().Info(ctx.Request().GetBody())
		r, err := ctx.Request().GetBody()
		if err != nil {
			return err
		}
		return r
	case "16":
		url := "http://192.168.5.108:8081/order/request/1"
		client, err := components.Def.HTTP().GetClient()
		if err != nil {
			return fmt.Errorf("GetClient():%+v", err)
		}
		status, content, err := client.Get(url)
		return types.XMap{
			"status":  status,
			"content": content,
			"err":     err,
		}
	case "17":

		sqlxxxx := `insert into sup_system_logs
		   (log_id, object_name, error_code, error_desc, create_time)
		   values
		   (seq_sys_log.nextval, 'test', 100, @err_desc, sysdate)`

		params := map[string]interface{}{
			"err_desc": time.Now().Format("20060102150405"),
		}

		dbObj := components.Def.DB().GetRegularDB()

		/*
			data, q, a, err := dbObj.Scalar("select 1 from dual", nil)
			ctx.Log().Info("dbObj.Scalar:", data, q, a, err)

				rows, q, a, err := dbObj.Query("select 1 r from dual", nil)
				ctx.Log().Info("dbObj.Query:", rows, q, a, err)
		*/

		effRow, q, a, err := dbObj.Execute(sqlxxxx, params)
		ctx.Log().Info("dbObj.Execute:", effRow, q, a, err)
		/*
			lastID, effRow, q, a, err := dbObj.Executes(sqlxxxx, params)
			ctx.Log().Info("dbObj.Executes:", lastID, effRow, q, a, err)

			trans, err := dbObj.Begin()

			data, q, a, err = trans.Scalar("select 1 from dual", nil)
			ctx.Log().Info("trans.Scalar:", data, q, a, err)

			rows, q, a, err = trans.Query("select 1 r from dual", nil)
			ctx.Log().Info("trans.Query:", rows, q, a, err)

			effRow, q, a, err = trans.Execute(sqlxxxx, params)
			ctx.Log().Info("trans.Execute:", effRow, q, a, err)

			lastID, effRow, q, a, err = trans.Executes(sqlxxxx, params)
			ctx.Log().Info("trans.Executes:", lastID, effRow, q, a, err)

			trans.Commit()
		*/
		return "db.test"

	case "18":
		cacheObj := components.Def.Cache().GetRegularCache()
		key := "hydra:newcache:apm"
		value := "10"
		keys := []string{key}
		delta := int64(2)
		expiresAt := 120

		err := cacheObj.Add(key, value, expiresAt)

		ctx.Log().Info("cacheObj.Add:", err)

		r, err := cacheObj.Get(key)
		ctx.Log().Info("cacheObj.Get:", r, err)

		v, err := cacheObj.Decrement(key, delta)

		ctx.Log().Info("cacheObj.Decrement:", v, err)
		v, err = cacheObj.Increment(key, delta)

		ctx.Log().Info("cacheObj.Increment:", v, err)
		rs, err := cacheObj.Gets(keys...)

		ctx.Log().Info("cacheObj.Gets:", rs, err)

		err = cacheObj.Set(key, value, expiresAt)

		ctx.Log().Info("cacheObj.Set:", err)
		err = cacheObj.Delay(key, expiresAt)

		ctx.Log().Info("cacheObj.Delay:", err)
		exists := cacheObj.Exists(key)

		ctx.Log().Info("cacheObj.Exists.1:", exists)
		err = cacheObj.Delete(key)

		ctx.Log().Info("cacheObj.Delete:", err)

		exists = cacheObj.Exists(key)

		ctx.Log().Info("cacheObj.Exists.2:", exists)
		return "cache.test"
	default:
		return hydra.G.PlatName
	}
}
