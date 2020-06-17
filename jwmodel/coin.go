package jwmodel

import (
	"encoding/json"
	"github.com/Alex-zs/oucjwparser/util"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"strconv"
)

type CoinStatus struct {
	Total float32 `json:"point_total"`	// 全部选课币
	Used float32 `json:"point_used"`		// 已使用选课比
	Rest float32 `json:"point_canused"`	// 剩余选课币
}

// 获取指定学年学期的选课币使用情况
func (session *JwSession) GetCoinStatus(year, semester int) (*CoinStatus, error){
	var coinStatusError error = nil
	coinStatus := new(CoinStatus)
	session.Do(func(req *fasthttp.Request, resp *fasthttp.Response) {
		// 设置请求参数
		req.SetRequestURI(Coin)
		req.Header.SetMethod(fasthttp.MethodPost)
		args := fasthttp.AcquireArgs()
		defer fasthttp.ReleaseArgs(args)
		args.Add("xh", session.UserCode)
		args.Add("xn", strconv.Itoa(year))
		args.Add("xq_m", strconv.Itoa(semester))
		args.WriteTo(req.BodyWriter())

		if err := fasthttp.DoTimeout(req, resp, util.DefaultHttpTimeout); err != nil {
			util.Log("请求选课比状态失败")
			coinStatusError = err
			return
		}

		resultStr := gjson.GetBytes(resp.Body(), "result").Str
		if err := json.Unmarshal([]byte(resultStr), coinStatus); err != nil {
			coinStatusError = err
		}
	})
	return coinStatus, coinStatusError
}
