package conf

import (
	"strings"
)

type HttpServerConf struct {
	Address           string `json:"address" valid:"dialstring"`
	Status            string `json:"status" valid:"in(start|stop)"`
	Engines           string `json:"engines"`
	ReadTimeout       int    `json:"readTimeout"`
	WriteTimeout      int    `json:"writeTimeout"`
	ReadHeaderTimeout int    `json:"readHeaderTimeout"`
	Name              string
	metadata          metadata
}

func (s *HttpServerConf) GetMetadata(key string) interface{} {
	return s.metadata.Get(key)
}
func (s *HttpServerConf) SetMetadata(key string, v interface{}) {
	s.metadata.Set(key, v)
}

type Metric struct {
	Host     string `json:"host" valid:"requrl,required"`
	DataBase string `json:"dataBase" valid:"ascii,required"`
	Cron     string `json:"cron" valid:"ascii,required"`
	UserName string `json:"userName" valid:"ascii"`
	Password string `json:"password" valid:"ascii"`
	Disable  bool   `json:"disable"`
}
type Authes map[string]*Auth

type Auth struct {
	Name       string   `json:"name" valid:"ascii,required"`
	ExpireAt   int64    `json:"expireAt" valid:"required"`
	Mode       string   `json:"mode" valid:"in(HS256|HS384|HS512|RS256|ES256|ES384|ES512|RS384|RS512|PS256|PS384|PS512),required"`
	Source     string   `json:"source" valid:"in(header|cookie|HEADER|COOKIE|H)"`
	Secret     string   `json:"secret" valid:"ascii,required"`
	Exclude    []string `json:"exclude"`
	FailedCode string   `json:"failed-code" valid:"numeric,range(400|999)"`
	Redirect   string   `json:"redirect" valid:"ascii"`
	Domain     string   `json:"domain" valid:"ascii"`
	Disable    bool     `json:"disable"`
}

var cacheExcludeSvcs = map[string]bool{}

//IsExcluded 是否是排除验证的服务
func (a *Auth) IsExcluded(service string) bool {

	service = strings.ToLower(service)
	if v, e := cacheExcludeSvcs[service]; e && v {
		return true
	}

	sparties := strings.Split(service, "/")
	//排除指定请求
	for _, u := range a.Exclude {
		//完全匹配
		if strings.EqualFold(u, service) {
			cacheExcludeSvcs[service] = true
			return true
		}
		//分段模糊
		uparties := strings.Split(u, "/")
		//取较少的数组长度
		uc := len(uparties)
		sc := len(sparties)
		/*
			处理模式：
			1. /a/b/ *
			2. /a/ **
			3. /a/ * /d
		**/
		if uc != sc && !strings.HasSuffix(u, "**") {
			continue
		}
		if uc > sc {
			continue
		}
		isMatch := true
		for i := 0; i < uc; i++ {
			if uparties[i] == "**" {
				cacheExcludeSvcs[service] = true
				return true
			}
			if uparties[i] == "*" {
				for j := i + 1; j < uc; j++ {
					if uparties[j] != sparties[j] {
						isMatch = false
						break
					}
				}
				if !isMatch {
					break
				}
				cacheExcludeSvcs[service] = true
				return true
			}
			if uparties[i] != sparties[i] {
				break
			}
		}

	}
	return false
}

type Routers struct {
	Setting map[string]string `json:"args"`
	Routers []*Router         `json:"routers"`
}
type Router struct {
	Name    string            `json:"name" valid:"ascii,required"`
	Action  []string          `json:"action" valid:"uppercase,in(GET|POST|PUT|DELETE|HEAD|TRACE|OPTIONS)"`
	Engine  string            `json:"engine" valid:"ascii,uppercase,in(*|RPC)"`
	Service string            `json:"service" valid:"ascii,required"`
	Setting map[string]string `json:"args"`
	Disable bool              `json:"disable"`
	Handler interface{}
}

type View struct {
	Path    string `json:"path" valid:"ascii,required"`
	Left    string `json:"left" valid:"ascii"`
	Right   string `json:"right" valid:"ascii"`
	Files   []string
	Disable bool `json:"disable"`
}

type CircuitBreaker struct {
	ForceBreak      bool       `json:"force-break"`
	Disable         bool       `json:"disable"`
	SwitchWindow    int        `json:"swith-window"`
	CircuitBreakers []*Breaker `json:"circuit-breakers"`
}
type Breaker struct {
	URL              string `json:"url" valid:"ascii,required"`
	RequestPerSecond int    `json:"request-per-second"`
	FailedPercent    int    `json:"failed-request"`
	RejectPerSecond  int    `json:"reject-per-second"`
	Disable          bool   `json:"disable"`
}
type Static struct {
	Dir       string   `json:"dir" valid:"ascii"`
	Prefix    string   `json:"prefix" valid:"ascii"`
	Exts      []string `json:"exts" valid:"ascii"`
	Exclude   []string `json:"exclude" valid:"ascii"`
	FirstPage string   `json:"first-page" valid:"ascii"`
	Rewriters []string `json:"rewriters" valid:"ascii"`
	Disable   bool     `json:"disable"`
}
type Tasks struct {
	Setting map[string]string `json:"args"`
	Tasks   []*Task           `json:"tasks"`
}
type Task struct {
	Name    string                 `json:"name" valid:"ascii"`
	Cron    string                 `json:"cron" valid:"ascii,required"`
	Input   map[string]interface{} `json:"input,omitempty"`
	Engine  string                 `json:"engine,omitempty"  valid:"ascii,uppercase,in(*|RPC)"`
	Service string                 `json:"service"  valid:"ascii,required"`
	Setting map[string]string      `json:"args"`
	Next    string                 `json:"next"`
	Last    string                 `json:"last"`
	Handler interface{}            `json:"handler,omitempty"`
	Disable bool                   `json:"disable"`
}

type Package struct {
	URL     string `json:"url" valid:"requrl,required"`
	Version string `json:"version" valid:"ascii,required"`
	CRC32   uint32 `json:"crc32" valid:"required"`
}

type Headers map[string]string
type Hosts []string

type Queues struct {
	Setting map[string]string `json:"args"`
	Queues  []*Queue          `json:"queues"`
}
type Server struct {
	Proto string `json:"proto" valid:"ascii,required"`
}
type Queue struct {
	Name        string            `json:"name" valid:"ascii"`
	Queue       string            `json:"queue" valid:"ascii,required"`
	Engine      string            `json:"engine,omitempty"  valid:"ascii,uppercase,in(*|RPC)"`
	Service     string            `json:"service" valid:"ascii,required"`
	Setting     map[string]string `json:"args"`
	Concurrency int               `json:"concurrency"`
	Disable     bool              `json:"disable"`
	Handler     interface{}
}
