package mcp

import (
	"net/http"
	"net/url"
	"strings"

	xjwt "github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/services"
)

// NewOAuthAuthorizeHandler returns the shared MCP OAuth authorize endpoint.
func NewOAuthAuthorizeHandler(cfg *services.MCPConf) func(context.IContext) interface{} {
	return func(ctx context.IContext) interface{} {
		if cfg == nil || cfg.OAuth == nil || cfg.OAuth.LoginURL == "" {
			return ctx.Response().JSON(http.StatusNotExtended, map[string]interface{}{
				"error": "mcp oauth login url is not configured",
			})
		}
		if data, ok := tryJWTAuth(ctx); ok {
			ctx.User().Auth().Request(data)
			if returnTo := safeReturnTo(ctx); returnTo != "" {
				ctx.Response().Redirect(http.StatusFound, returnTo)
				return nil
			}
			return map[string]interface{}{
				"authorized": true,
				"scopes":     cfg.OAuth.Scopes,
			}
		}
		ctx.Response().Redirect(http.StatusFound, loginRedirectURL(cfg.OAuth.LoginURL, ctx))
		return nil
	}
}

func tryJWTAuth(ctx context.IContext) (interface{}, bool) {
	jwtConf, err := ctx.APPConf().GetJWTConf()
	if err != nil || jwtConf.Disable {
		return nil, false
	}
	data, err := jwtConf.CheckJWT(tokenFromRequest(ctx, jwtConf))
	return data, err == nil
}

func tokenFromRequest(ctx context.IContext, jwtConf *xjwt.JWTAuth) string {
	// authorize 可能由浏览器导航触发(登录成功后跳转), 浏览器导航只能携带 cookie,
	// 无法携带 Authorization header。因此按全局 source 读取失败时, 兜底读取另一个 source,
	// 使 OAuth 浏览器跳转场景(走 cookie)与业务 API(走 header)能各自命中。
	primary, fallback := sourcePair(jwtConf.Source)
	if token := readTokenBySource(ctx, jwtConf, primary); token != "" {
		return token
	}
	return readTokenBySource(ctx, jwtConf, fallback)
}

// sourcePair 返回 (主来源, 兜底来源): header 与 cookie 互为兜底。
func sourcePair(source string) (primary, fallback string) {
	switch strings.ToUpper(source) {
	case xjwt.SourceHeader, xjwt.SourceHeaderShort:
		return xjwt.SourceHeader, xjwt.SourceCookie
	default:
		return xjwt.SourceCookie, xjwt.SourceHeader
	}
}

// readTokenBySource 按指定来源(header/cookie)读取 token。
func readTokenBySource(ctx context.IContext, jwtConf *xjwt.JWTAuth, source string) string {
	switch strings.ToUpper(source) {
	case xjwt.SourceHeader, xjwt.SourceHeaderShort:
		return ctx.Request().Headers().GetString(xjwt.AuthorizationHeader)
	default:
		return ctx.Request().Cookies().GetString(jwtConf.Name)
	}
}

func loginRedirectURL(loginURL string, ctx context.IContext) string {
	u, err := url.Parse(loginURL)
	if err != nil {
		return loginURL
	}
	q := u.Query()
	q.Set("return_to", currentRequestURL(ctx))
	u.RawQuery = q.Encode()
	return u.String()
}

func safeReturnTo(ctx context.IContext) string {
	reqURL := ctx.Request().Path().GetURL()
	if reqURL == nil {
		return ""
	}
	raw := reqURL.Query().Get("return_to")
	if raw == "" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	if !u.IsAbs() {
		if strings.HasPrefix(raw, "/") {
			return raw
		}
		return ""
	}
	if u.Host == requestHost(ctx, reqURL) {
		return u.String()
	}
	return ""
}

func currentRequestURL(ctx context.IContext) string {
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

func requestOrigin(ctx context.IContext, reqURL *url.URL) string {
	scheme := firstHeaderValue(ctx.Request().Headers().GetString("X-Forwarded-Proto"))
	if scheme == "" && reqURL != nil {
		scheme = reqURL.Scheme
	}
	if scheme == "" {
		scheme = "http"
	}
	host := requestHost(ctx, reqURL)
	if host == "" {
		return ""
	}
	return scheme + "://" + host
}

func requestHost(ctx context.IContext, reqURL *url.URL) string {
	host := firstHeaderValue(ctx.Request().Headers().GetString("X-Forwarded-Host"))
	if host == "" {
		host = ctx.Request().Headers().GetString("Host")
	}
	if host == "" && reqURL != nil {
		host = reqURL.Host
	}
	return host
}

func firstHeaderValue(value string) string {
	if i := strings.Index(value, ","); i >= 0 {
		return strings.TrimSpace(value[:i])
	}
	return strings.TrimSpace(value)
}
