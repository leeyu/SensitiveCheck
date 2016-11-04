package sensitivecheck

import (
	"net/http"
)

type Context struct {
	Request     *http.Request
	Response    http.ResponseWriter
	QueryParams map[string]string
	RouteParams map[string]string
	BodyParams  map[string]string
}

func NewContext(params map[string]string, res http.ResponseWriter, req *http.Request) *Context {
	ctx := Context{}
	req.ParseForm()
	ctx.Request = req
	ctx.Response = res
	ctx.RouteParams = params
	return &ctx
}

func (ctx *Context) GetParam(key string, default_val interface{}) string {
	val, exists := ctx.RouteParams[key]
	if exists {
		return val
	}

	val, exists = ctx.QueryParams[key]
	if exists {
		return val
	}

	val, exists = ctx.BodyParams[key]
	if exists {
		return val
	}

	return default_val.(string)
}
