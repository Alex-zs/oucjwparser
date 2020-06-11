package util

import (
	"github.com/valyala/fasthttp"
	"time"
)

const DefaultTimeout = 10 * time.Second

// 请求复用
func SimpleDo(f func(req *fasthttp.Request, resp *fasthttp.Response))  {
	// 请求、响应对象复用
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	// 退出时释放对象
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()
	f(req, resp)
}

