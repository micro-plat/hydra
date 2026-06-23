// Package mcp 为 Hydra 提供 MCP(Model Context Protocol) 可选特性支持。
//
// 它通过 services.RegisterTypedHandlerHook 注入 typed Handle 识别能力，把
// func(ctx context.IContext, req)(resp, error) 签名包装为 context.Handler，
// 供 HTTP 自动参数绑定与 MCP schema 反射使用；并通过 services.RegisterMCPHook
// 注入 tool 注册能力，供 hydra.S.Micro(...).WithMCP() / hydra.S.MCP(...) 登记 tool。
//
// 本包为可选特性包（位于仓库根，与 context/services/conf 等核心包平级），只依赖
// 核心层、不依赖 hydra/hydra/servers 应用层。未被 import 时框架行为与现状完全一致。
package mcp

import (
	"reflect"
	"strings"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/services"
)

// iContextType context.IContext 接口的反射类型，用于判断第1参是否实现该接口。
var iContextType = reflect.TypeOf((*context.IContext)(nil)).Elem()

// httpMethods 镜像 conf/server/router.Methods，用于判断对象方法名是否为 RESTful 动词。
var httpMethods = map[string]bool{
	"GET": true, "POST": true, "PUT": true, "DELETE": true, "OPTIONS": true, "HEAD": true,
}

// resolved 一个 typed Handle 的解析结果：源 HTTP 路径、包装后的 handler、输入/输出类型、
// 方法名（对象服务用于查 mcp.Service 声明的描述）。
type resolved struct {
	sourcePath string
	handler    context.Handler
	reqType    reflect.Type
	respType   reflect.Type
	methodName string // 对象服务的 Handle 方法名（如 PostHandle）；函数形式为 ""
}

// recognizeTyped 判定 typed Handle 签名 func(ctx context.IContext, req)(resp, error)，
// 返回 req/resp 类型。非该签名返回 ok=false。
func recognizeTyped(i interface{}) (reqType, respType reflect.Type, ok bool) {
	typ := reflect.TypeOf(i)
	if typ == nil || typ.Kind() != reflect.Func {
		return nil, nil, false
	}
	// typed Handle: 两入参(ctx, req)、两返回值(resp, error)
	if typ.NumIn() != 2 || typ.NumOut() != 2 {
		return nil, nil, false
	}
	// 第1参须实现 context.IContext
	if !typ.In(0).Implements(iContextType) {
		return nil, nil, false
	}
	return typ.In(1), typ.Out(0), true
}

// wrapValue 按 req 类型把 typed 函数值包装为 context.Handler：Bind req、调原方法、
// 返回值经框架 WriteAny（error 走 400 分支，与旧签名语义一致）。
func wrapValue(fn reflect.Value, reqType reflect.Type) context.Handler {
	return func(ctx context.IContext) interface{} {
		ptr := newPtr(reqType)
		if err := ctx.Request().Bind(ptr); err != nil {
			return err
		}
		var reqVal reflect.Value
		if reqType.Kind() == reflect.Ptr {
			reqVal = reflect.ValueOf(ptr)
		} else {
			reqVal = reflect.ValueOf(ptr).Elem()
		}
		out := fn.Call([]reflect.Value{reflect.ValueOf(ctx), reqVal})
		if err, ok := out[1].Interface().(error); ok && err != nil {
			return err
		}
		return out[0].Interface()
	}
}

// wrapTyped services.RegisterTypedHandlerHook 的注入目标：识别 typed 函数并包装为 context.Handler。
// 严格判定 NumIn()==2 && NumOut()==2，不会误伤 Handling/Handled/Fallback 钩子（其返回值个数/参数不匹配）。
func wrapTyped(i interface{}) (context.Handler, bool) {
	reqType, _, ok := recognizeTyped(i)
	if !ok {
		return nil, false
	}
	return wrapValue(reflect.ValueOf(i), reqType), true
}

// resolve 解析 h（typed 函数 或 含 typed *Handle 方法的对象），basePath 为注册路径。
// 返回每个 typed Handle 对应的 resolved；无 typed Handle 时 ok=false（旧签名不登记 MCP tool）。
func resolve(h interface{}, basePath string) (items []resolved, ok bool) {
	if h == nil {
		return nil, false
	}
	// 函数形式
	if reqType, respType, typed := recognizeTyped(h); typed {
		return []resolved{{
			sourcePath: basePath,
			handler:    wrapValue(reflect.ValueOf(h), reqType),
			reqType:    reqType,
			respType:   respType,
		}}, true
	}
	// 对象形式：遍历 *Handle 方法（镜像 services.reflectHandle 的方法发现）
	typ := reflect.TypeOf(h)
	if typ.Kind() == reflect.Func {
		return nil, false
	}
	val := reflect.ValueOf(h)
	for i := 0; i < typ.NumMethod(); i++ {
		mName := typ.Method(i).Name
		if !strings.HasSuffix(mName, handleSuffix) {
			continue
		}
		method := val.MethodByName(mName)
		reqType, respType, typed := recognizeTyped(method.Interface())
		if !typed {
			continue // 非 typed Handle（旧签名钩子）不作为 MCP tool
		}
		endName := camelToPath(strings.TrimSuffix(mName, handleSuffix))
		items = append(items, resolved{
			sourcePath: resolveSourcePath(basePath, endName),
			handler:    wrapValue(method, reqType),
			reqType:    reqType,
			respType:   respType,
			methodName: mName,
		})
	}
	return items, len(items) > 0
}

// resolveSourcePath 镜像 services.UnitGroup.getPaths 的源路径规则，计算 tool 的源 HTTP 路径。
// endName 为对象方法去掉 Handle 后缀经 camelToPath 的值（函数形式为 ""）。
func resolveSourcePath(basePath, endName string) string {
	if endName == "" {
		return basePath // 函数形式：路径即源路径
	}
	if httpMethods[strings.ToUpper(endName)] {
		return basePath // RESTful 动词方法：路径不变，方法名决定 HTTP 动词
	}
	return joinPath(basePath, endName) // 非动词：拼接子路径
}

// newPtr 按 req 类型创建可写的指针实例。
// req 为 struct 时返回 *struct；req 为 *struct 时返回指向同元素类型的指针。
func newPtr(reqType reflect.Type) interface{} {
	if reqType.Kind() == reflect.Ptr {
		return reflect.New(reqType.Elem()).Interface()
	}
	return reflect.New(reqType).Interface()
}

func init() {
	services.RegisterTypedHandlerHook(wrapTyped)
}
