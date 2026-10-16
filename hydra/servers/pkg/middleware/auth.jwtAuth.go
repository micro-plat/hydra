package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	xjwt "github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/errs"
)

// JwtAuth jwt
func JwtAuth() Handler {

	return func(ctx IMiddleContext) {

		//1. 获取jwt配置
		jwtAuth, err := ctx.APPConf().GetJWTConf()
		if err != nil {
			ctx.Response().Abort(xjwt.JWTStatusConfError, err)
			return
		}
		if jwtAuth.Disable {
			ctx.Next()
			return
		}

		//2.检查jwt是否有效
		ctx.Response().AddSpecial("jwt")
		_, jwtErr := checkJWT(ctx, jwtAuth)

		//3.检查是否需要跳过请求
		if ok, _ := jwtAuth.Match(ctx.Request().Path().GetRequestPath()); ok {
			ctx.Next()
			return
		}

		//4. 验证jwt
		if jwtErr == nil {
			ctx.Next()
			return
		}

		//5.jwt验证失败后返回错误
		ctx.Log().Error(jwtErr)
		authenticateURL := writeJWTAuthenticateHeader(ctx, jwtAuth)
		if authenticateURL != "" && isBrowserNavigation(ctx) {
			ctx.Response().Redirect(http.StatusFound, authenticateURL)
			return
		}
		if jwtAuth.AuthURL != "" {
			ctx.Response().Header(header.XLocation, ctx.Request().Headers().Translate(jwtAuth.AuthURL))
			ctx.Response().Abort(xjwt.JWTStatusTokenError)
			return
		}
		ctx.Response().Abort(errs.GetCode(jwtErr, xjwt.JWTStatusTokenError), errors.New("jwt验证串错误，禁止访问"))
	}
}

func writeJWTAuthenticateHeader(ctx IMiddleContext, jwtAuth *xjwt.JWTAuth) string {
	path := ctx.Request().Path().GetRequestPath()
	challengeURL := jwtAuth.GetAuthenticateURL(path)
	if challengeURL == "" {
		return ""
	}
	// 仅显式配置了路径级认证URL的路径才附加 return_to；
	// 未配置的路径（如 API）登录后跳默认页，不回跳到当前请求URL。
	authenticateURL := challengeURL
	if jwtAuth.HasPathAuthenticateURL(path) {
		authenticateURL = appendReturnTo(challengeURL, ctx)
	}
	ctx.Response().Header("WWW-Authenticate", fmt.Sprintf(`Bearer location="%s"`, authenticateURL))
	return authenticateURL
}

func isBrowserNavigation(ctx IMiddleContext) bool {
	if strings.ToUpper(ctx.Request().Path().GetMethod()) != http.MethodGet {
		return false
	}
	accept := ctx.Request().Headers().GetString("Accept")
	secFetchMode := ctx.Request().Headers().GetString("Sec-Fetch-Mode")
	return strings.Contains(accept, "text/html") || strings.EqualFold(secFetchMode, "navigate")
}

func appendReturnTo(challengeURL string, ctx IMiddleContext) string {
	u, err := url.Parse(challengeURL)
	if err != nil {
		return challengeURL
	}
	q := u.Query()
	if q.Get("return_to") == "" {
		q.Set("return_to", currentRequestURL(ctx))
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func currentRequestURL(ctx IMiddleContext) string {
	reqURL := ctx.Request().Path().GetURL()
	path := ctx.Request().Path().GetRequestPath()
	if path == "" && reqURL != nil {
		path = reqURL.EscapedPath()
	}
	if path == "" {
		path = "/"
	}
	if reqURL != nil && reqURL.RawQuery != "" {
		path += "?" + reqURL.RawQuery
	}
	if origin := requestOrigin(ctx, reqURL); origin != "" {
		return origin + path
	}
	return path
}

func requestOrigin(ctx IMiddleContext, reqURL *url.URL) string {
	scheme := firstHeaderValue(ctx.Request().Headers().GetString("X-Forwarded-Proto"))
	if scheme == "" && reqURL != nil {
		scheme = reqURL.Scheme
	}
	if scheme == "" {
		scheme = "http"
	}
	host := firstHeaderValue(ctx.Request().Headers().GetString("X-Forwarded-Host"))
	if host == "" {
		host = ctx.Request().Headers().GetString("Host")
	}
	if host == "" && reqURL != nil {
		host = reqURL.Host
	}
	if host == "" {
		return ""
	}
	return scheme + "://" + host
}

func firstHeaderValue(value string) string {
	if i := strings.Index(value, ","); i >= 0 {
		return strings.TrimSpace(value[:i])
	}
	return strings.TrimSpace(value)
}

// CheckJWT 检查jwk参数是否合法
func checkJWT(ctx context.IContext, j *xjwt.JWTAuth) (data interface{}, err error) {

	//1. 从请求中获取jwt信息
	token := getToken(ctx, j)
	data, err = j.CheckJWT(token)
	if err != nil {
		return nil, err
	}

	//保存到Context中
	ctx.User().Auth().Request(data)
	return data, nil
}

// getToken 从请求头或cookie中获取cookie
func getToken(ctx context.IContext, jwt *xjwt.JWTAuth) string {
	switch strings.ToUpper(jwt.Source) {
	case xjwt.SourceHeader, xjwt.SourceHeaderShort:
		return ctx.Request().Headers().GetString(xjwt.AuthorizationHeader)
	default:
		cookie := ctx.Request().Cookies().GetString(jwt.Name)
		return cookie
	}
}
