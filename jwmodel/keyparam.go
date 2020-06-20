package jwmodel

import (
	"github.com/Alex-zs/oucjwparser/util"
	"github.com/valyala/fasthttp"
	"strings"
)

// 教务系统的加密参数
// 部分请求需要加密参数
type KeyParam struct {
	Key  string
	Time string
}

// 获取加密参数
func (session *JwSession) GetKeyParam() (*KeyParam, error) {
	var keyParamError error = nil
	keyParam := new(KeyParam)
	session.Do(func(req *fasthttp.Request, resp *fasthttp.Response) {
		req.SetRequestURI(KeyTimeParam)
		req.Header.SetMethod(fasthttp.MethodGet)
		if err := fasthttp.DoTimeout(req, resp, util.DefaultHttpTimeout); err != nil {
			keyParamError = err
			return
		}
		data := string(resp.Body())
		splice := strings.Split(data, "\n")
		if len(splice) < 4 || len(splice[2]) < 34 || len(splice[3]) < 35 {
			keyParamError = &util.JwError{Msg: "获取加密参数错误"}
			return
		}
		keyParam.Key = splice[2][15:34]
		keyParam.Time = splice[3][16:35]
	})
	return keyParam, keyParamError
}
