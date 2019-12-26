package http

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/servers"
	"github.com/micro-plat/hydra/servers/http/middleware"
	"github.com/micro-plat/lib4go/archiver"
	"github.com/micro-plat/lib4go/types"
)

//waitRemoveDir 等待移除的静态文件
var waitRemoveDir = make([]string, 0, 1)

type ISetMetric interface {
	SetMetric(*conf.Metric) error
}

//SetMetric 设置metric
func SetMetric(set ISetMetric, cnf conf.IServerConf) (enable bool, err error) {
	//设置静态文件路由
	var metric conf.Metric
	_, err = cnf.GetSubObject("metric", &metric)

	if err != nil && err != conf.ErrNoSetting {
		return false, err
	}
	if err == conf.ErrNoSetting {
		metric.Disable = true
	} else {
		if b, err := govalidator.ValidateStruct(&metric); !b {
			err = fmt.Errorf("metric配置有误:%v", err)
			return false, err
		}
	}
	err = set.SetMetric(&metric)
	return !metric.Disable && err == nil, err
}

type ISetStatic interface {
	SetStatic(static *conf.Static) error
}

//SetStatic 设置static
func SetStatic(set ISetStatic, cnf conf.IServerConf) (enable bool, err error) {
	//设置静态文件路由
	var static conf.Static
	_, err = cnf.GetSubObject("static", &static)
	if err != nil && err != conf.ErrNoSetting {
		return false, err
	}
	if err != conf.ErrNoSetting {
		if b, err := govalidator.ValidateStruct(&static); !b {
			err = fmt.Errorf("static配置有误:%v", err)
			return false, err
		}
	}
	if static.Dir == "" {
		static.Dir = "./static"
	}
	if static.FirstPage == "" {
		static.FirstPage = "index.html"
	}
	static.Exts = append(static.Exts, ".txt", ".jpg", ".png", ".gif", ".ico", ".html", ".htm", ".js", ".css", ".map", ".ttf", ".woff", ".woff2", ".woff2")
	static.Rewriters = append(static.Rewriters, "/", "index.htm", "default.html")
	static.Exclude = append(static.Exclude, "/views/", ".exe", ".so")
	static.Dir, err = unarchive(static.Dir, static.Archive) //处理归档文件
	if err != nil {
		return false, err
	}
	err = set.SetStatic(&static)
	return !static.Disable && err == nil, err
}

//ISetRouterHandler 设置路由列表
type ISetRouterHandler interface {
	SetRouters([]*conf.Router) error
}

func getDefaultRouters(services map[string][]string) []*conf.Router {
	routers := conf.Routers{}

	if len(services) == 0 {
		routers.Routers = make([]*conf.Router, 0, 1)
		routers.Routers = append(routers.Routers, &conf.Router{Action: []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}, Name: "/*name", Service: "/@name", Engine: "*"})
		return routers.Routers
	}
	routers.Routers = make([]*conf.Router, 0, len(services))
	for name, actions := range services {
		router := &conf.Router{
			Action:  actions, //[]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
			Name:    name,
			Service: name,
			Engine:  "*",
		}
		router.Action = append(router.Action, "OPTIONS")
		routers.Routers = append(routers.Routers, router)
	}
	return routers.Routers
}

//SetHttpRouters 设置路由
func SetHttpRouters(engine servers.IRegistryEngine, set ISetRouterHandler, cnf conf.IServerConf) (enable bool, err error) {

	var routers conf.Routers
	if _, err = cnf.GetSubObject("router", &routers); err == conf.ErrNoSetting || len(routers.Routers) == 0 {
		routers.Routers = getDefaultRouters(engine.GetServices()) //添加默认路由
	}
	if err != nil && err != conf.ErrNoSetting {
		err = fmt.Errorf("路由:%v", err)
		return false, err
	}
	if b, err := govalidator.ValidateStruct(&routers); !b {
		err = fmt.Errorf("router配置有误:%v", err)
		return false, err
	}

	//处理RPC代理服务
	nRouters := make([]*conf.Router, 0, len(routers.RPCS)+len(routers.Routers))
	for _, proxy := range routers.RPCS {
		if len(proxy.Action) == 0 {
			proxy.Action = []string{"GET", "POST"}
		}
		proxy.Engine = "rpc"
		nRouters = append(nRouters, proxy)
	}

	//处理路由
	for _, router := range routers.Routers {
		if len(router.Action) == 0 {
			router.Action = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}
		}
		router.Engine = types.GetString(router.Engine, "*")
		nRouters = append(nRouters, router)
	}

	//构建路由处理函数
	for _, router := range nRouters {
		if router.Setting == nil {
			router.Setting = make(map[string]string)
		}
		for k, v := range routers.Setting {
			if _, ok := router.Setting[k]; !ok {
				router.Setting[k] = v
			}
		}
		router.Handler = middleware.ContextHandler(engine, router.Name, router.Engine, router.Service, router.Setting)

	}
	err = set.SetRouters(nRouters)
	return len(nRouters) > 0 && err == nil, err
}

//---------------------------------------------------------------------------
//-------------------------------view---------------------------------------
//---------------------------------------------------------------------------

//ISetView 设置view
type ISetView interface {
	SetView(*conf.View) error
}

//SetView 设置view
func SetView(set ISetView, cnf conf.IServerConf) (enable bool, err error) {
	//设置jwt安全认证参数
	var view conf.View
	_, err = cnf.GetSubObject("view", &view)
	if err != nil && err != conf.ErrNoSetting {
		return false, err
	}
	if err == conf.ErrNoSetting {
		view.Disable = true
	} else {
		if b, err := govalidator.ValidateStruct(&view); !b {
			err = fmt.Errorf("view配置有误:%v", err)
			return false, err
		}
	}
	err = set.SetView(&view)
	return err == nil && !view.Disable, err
}

//ISetCircuitBreaker 设置CircuitBreaker
type ISetCircuitBreaker interface {
	CloseCircuitBreaker() error
	SetCircuitBreaker(*conf.CircuitBreaker) error
}

//SetCircuitBreaker 设置熔断配置
func SetCircuitBreaker(set ISetCircuitBreaker, cnf conf.IServerConf) (enable bool, err error) {
	//设置CircuitBreaker
	var breaker conf.CircuitBreaker
	if _, err = cnf.GetSubObject("circuit", &breaker); err == conf.ErrNoSetting || breaker.Disable {
		return false, set.CloseCircuitBreaker()
	}
	if err != nil {
		return false, err
	}
	if b, err := govalidator.ValidateStruct(&breaker); !b {
		err = fmt.Errorf("circuit配置有误:%v", err)
		return false, err
	}
	err = set.SetCircuitBreaker(&breaker)
	return err == nil && !breaker.Disable, err
}

//---------------------------------------------------------------------------
//-------------------------------header---------------------------------------
//---------------------------------------------------------------------------

//ISetHeaderHandler 设置header
type ISetHeaderHandler interface {
	SetHeader(conf.Headers) error
}

//SetHeaders 设置header
func SetHeaders(set ISetHeaderHandler, cnf conf.IServerConf) (enable bool, err error) {
	//设置通用头信息
	var header conf.Headers
	_, err = cnf.GetSubObject("header", &header)
	if err != nil && err != conf.ErrNoSetting {
		err = fmt.Errorf("header配置有误:%v", err)
		return false, err
	}
	err = set.SetHeader(header)
	return len(header) > 0 && err == nil, err
}

//---------------------------------------------------------------------------
//-------------------------------ajax---------------------------------------
//---------------------------------------------------------------------------

//IAjaxRequest 设置ajax
type IAjaxRequest interface {
	SetAjaxRequest(bool) error
}

//SetAjaxRequest 设置ajax
func SetAjaxRequest(set IAjaxRequest, cnf conf.IServerConf) (enable bool, err error) {
	if enable = cnf.GetBool("onlyAllowAjaxRequest", false); !enable {
		return false, nil
	}
	err = set.SetAjaxRequest(enable)
	return enable && err == nil, err
}

//---------------------------------------------------------------------------
//-------------------------------host---------------------------------------
//---------------------------------------------------------------------------

//ISetHosts 设置hosts
type ISetHosts interface {
	SetHosts(conf.Hosts) error
}

//SetHosts 设置hosts
func SetHosts(set ISetHosts, cnf conf.IServerConf) (enable bool, err error) {
	var hosts conf.Hosts
	hosts = cnf.GetStrings("host")
	err = set.SetHosts(hosts)
	return len(hosts) > 0 && err == nil, err
}

//---------------------------------------------------------------------------
//-------------------------------response---------------------------------------
//---------------------------------------------------------------------------

//ISetResponse 设置response配置串
type ISetResponse interface {
	SetResponse(*conf.Response) error
}

//SetResponse 设置response配置串
func SetResponse(set ISetResponse, cnf conf.IServerConf) (enable bool, err error) {
	var response conf.Response
	_, err = cnf.GetSubObject("response", &response)
	if err != nil && err != conf.ErrNoSetting {
		return false, err
	}
	if err == conf.ErrNoSetting {
		return false, nil
	}
	err = set.SetResponse(&response)
	return err == nil, err
}

//---------------------------------------------------------------------------
//-------------------------------jwt---------------------------------------
//---------------------------------------------------------------------------

//ISetJwtAuth 设置jwt
type ISetJwtAuth interface {
	SetJWT(*conf.JWTAuth) error
}

//SetJWT 设置jwt
func SetJWT(set ISetJwtAuth, cnf conf.IServerConf) (enable bool, err error) {
	//设置jwt安全认证参数
	var auths conf.Authes
	if _, err := cnf.GetSubObject("auth", &auths); err != nil && err != conf.ErrNoSetting {
		err = fmt.Errorf("jwt配置有误:%v", err)
		return false, err
	}
	if auths.JWT != nil {
		if b, err := govalidator.ValidateStruct(auths.JWT); !b {
			err = fmt.Errorf("jwt配置有误:%v", err)
			return false, err
		}
		err = set.SetJWT(auths.JWT)
		return err == nil && !auths.JWT.Disable, err
	}
	return false, nil
}

func unarchive(dir string, path string) (string, error) {
	if path == "" {
		return dir, nil
	}
	archive := archiver.MatchingFormat(path)
	if archive == nil {
		return "", fmt.Errorf("指定的文件不是归档文件:%s", path)
	}
	tmpDir, err := ioutil.TempDir("", "hydra")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败:%v", err)
	}
	reader, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("无法打开文件:%s(%v)", path, err)
	}
	defer reader.Close()
	err = archive.Read(reader, tmpDir)
	if err != nil {
		return "", fmt.Errorf("读取归档文件失败:%v", err)
	}
	ndir := filepath.Join(tmpDir, dir)
	waitRemoveDir = append(waitRemoveDir, tmpDir)
	return ndir, nil
}

//---------------------------------------------------------------------------
//-------------------------------fixed-secret---------------------------------------
//---------------------------------------------------------------------------

//CheckFixedSecret 设置FixedSecret
func CheckFixedSecret(cnf conf.IServerConf) (enable bool, err error) {
	//设置fixedSecret安全认证参数
	var auths conf.Authes
	if _, err := cnf.GetSubObject("auth", &auths); err != nil && err != conf.ErrNoSetting {
		err = fmt.Errorf("fixed-secret配置有误:%v", err)
		return false, err
	}
	if auths.FixedScret != nil {
		if b, err := govalidator.ValidateStruct(auths.FixedScret); !b {
			err = fmt.Errorf("fixed-secret配置有误:%v", err)
			return false, err
		}
		return !auths.FixedScret.Disable, nil
	}
	return false, nil
}

//---------------------------------------------------------------------------
//-------------------------------remote-auth---------------------------------------
//---------------------------------------------------------------------------

//CheckRemoteAuth 检查是否设置remote-auth
func CheckRemoteAuth(cnf conf.IServerConf) (enable bool, err error) {
	//设置Remote安全认证参数
	var auths conf.Authes
	if _, err := cnf.GetSubObject("auth", &auths); err != nil && err != conf.ErrNoSetting {
		err = fmt.Errorf("remote-auth配置有误:%v", err)
		return false, err
	}
	count := 0
	for _, auth := range auths.RemotingServiceAuths {
		if b, err := govalidator.ValidateStruct(auth); !b {
			err = fmt.Errorf("remote-auth配置有误:%v", err)
			return false, err
		}
		if !auth.Disable {
			count++
		}
	}
	return count > 0, nil
}
