package dispatcher

import (
	"sync"
)

var default404Body = []byte("404 service not found")

type Dispatcher struct {
	RouterGroup
	trees            methodTrees
	pool             sync.Pool
	secureJsonPrefix string
}

func New() *Dispatcher {
	engine := &Dispatcher{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		trees:            make(methodTrees, 0, 9),
		secureJsonPrefix: "while(1);",
	}
	engine.RouterGroup.engine = engine
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	return engine
}
func (engine *Dispatcher) allocateContext() *Context {
	return &Context{engine: engine}
}
func (engine *Dispatcher) addRoute(method, path string, handlers HandlersChain) {
	root := engine.trees.get(method)
	if root == nil {
		root = new(node)
		engine.trees = append(engine.trees, methodTree{method: method, root: root})
	}
	root.addRoute(path, handlers)
}

func (engine *Dispatcher) Use(middleware ...HandlerFunc) IRoutes {
	engine.RouterGroup.Use(middleware...)
	return engine
}

// Routes returns a slice of registered routes, including some useful information, such as:
// the http method, path and the handler name.
func (engine *Dispatcher) Routes() (routes RoutesInfo) {
	for _, tree := range engine.trees {
		routes = iterate("", tree.method, routes, tree.root)
	}
	return routes
}

func iterate(path, method string, routes RoutesInfo, root *node) RoutesInfo {
	path += root.path
	if len(root.handlers) > 0 {
		routes = append(routes, RouteInfo{
			Method:  method,
			Path:    path,
			Handler: nameOfFunction(root.handlers.Last()),
		})
	}
	for _, child := range root.children {
		routes = iterate(path, method, routes, child)
	}
	return routes
}

//HandleRequest 处理外部请求
func (engine *Dispatcher) HandleRequest(r IRequest) (w *responseWriter, err error) {
	c := engine.pool.Get().(*Context)
	c.reset(r)
	c.writermem.reset()
	engine.handleRequest(c)
	if len(c.Errors) > 0 {
		err = c.Errors[0]
	}
	w = c.writermem.Copy()
	c.writermem.reset()
	c.reset(nil)
	engine.pool.Put(c)
	return
}

//Find 查找指定路由是否存在
func (engine *Dispatcher) Find(path string) bool {
	t := engine.trees
	for i, tl := 0, len(t); i < tl; i++ {
		root := t[i].root
		handlers, _, _ := root.getValue(path, make(Params, 0, 2), false)
		if handlers != nil {
			return true
		}
	}
	return false
}
func (engine *Dispatcher) handleRequest(c *Context) {
	httpMethod := c.Request.GetMethod()
	path := c.Request.GetService()
	unescape := false
	if httpMethod == "" {
		httpMethod = "GET"
	}
	// Find root of the tree for the given HTTP method
	t := engine.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method == httpMethod {
			root := t[i].root
			// Find route in tree
			handlers, params, _ := root.getValue(path, c.Params, unescape)
			if handlers != nil {
				c.handlers = handlers
				c.Params = params
				c.Next()
				c.writermem.WriteHeaderNow()
				return
			}
			break
		}
	}
	c.handlers = engine.RouterGroup.Handlers
	serveError(c, 404, default404Body)
}

var mimePlain = []string{"text/plain"}

func serveError(c *Context, code int, defaultMessage []byte) {
	c.writermem.status = code
	c.Next()
	if !c.writermem.Written() {
		if c.writermem.Status() == code {
			c.writermem.Header()["Content-Type"] = mimePlain
			c.Writer.Write(defaultMessage)
		} else {
			c.writermem.WriteHeaderNow()
		}
	}
}
