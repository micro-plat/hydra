package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/examples/apiserver/services/order"
	"github.com/micro-plat/hydra/hydra"
)

func main() {

	app := hydra.NewApp(
		hydra.WithPlatName("hydra-780"),
		hydra.WithSystemName("apiserver"),
		hydra.WithServerTypes("api"),
		hydra.WithDebug())
	app.Handling(func(c *context.Context) interface{} {
		fmt.Println("go id:", Goid())
		return nil
	})
	app.Micro("/order/query", order.NewQueryHandler)
	app.Micro("/order/bind", order.NewBindHandler)
	app.Start()
}
func Goid() int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic recover:panic info:%v", err)
		}
	}()

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
