package util

import (
	"github.com/valyala/fasthttp"
	"time"
)

// 默认http超时时间
const DefaultHttpTimeout = 10 * time.Second

// 简单的请求复用
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

