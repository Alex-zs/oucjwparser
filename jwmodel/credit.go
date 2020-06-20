package jwmodel

import (
	"bytes"
	"github.com/Alex-zs/oucjwparser/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
	"strconv"
)

// 毕业学分要求
type CreditRequire struct {
	CreditType string		// 学分类型，如 `专业知识/限选课`
	Credit float64			// 这类课程总计所需的学分
}

// 获取毕业要求的每类课程所需的学分
func (session *JwSession) GetCreditRequire() ([]CreditRequire, error) {
	var creditList []CreditRequire
	var creditError error
	session.Do(func(req *fasthttp.Request, resp *fasthttp.Response) {
		// 设置请求参数
		req.SetRequestURI(DataTable)
		req.Header.SetMethod(fasthttp.MethodPost)
		args := fasthttp.AcquireArgs()
		defer fasthttp.ReleaseArgs(args)
		args.Add("tableId", "6033")
		args.WriteTo(req.BodyWriter())
		// 请求
		if err := fasthttp.DoTimeout(req, resp, util.DefaultHttpTimeout); err != nil {
			util.Log("请求学分数据失败")
			creditError = err
			return
		}
		data := util.GBK2UTF8(resp.Body())
		creditList, creditError = parseHtmlCredit(data)
	})
	return creditList, creditError
}

func parseHtmlCredit(data []byte) ([]CreditRequire, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	var parseError error

	var creditList = make([]CreditRequire, 0)
	doc.Find("tbody").Find("tr").Each(func(i int, selection *goquery.Selection) {
		credit := CreditRequire{}
		selection.Find("td").Each(func(i int, selection *goquery.Selection) {
			switch i {
			case 1:
				credit.CreditType = selection.Text()
			case 2:
				credit.Credit, err = strconv.ParseFloat(selection.Text(), 64)
				if err != nil {
					parseError = err
					return
				}
			}
		})
		if credit.CreditType != "合计" {
			creditList = append(creditList, credit)
		}
	})

	// 根据通识教育/限选课计算正常学分
	multiple := 1.0
	for _, credit := range creditList {
		if credit.CreditType == "通识教育/限选课" {
			multiple = credit.Credit / 8.0
			break
		}
	}
	for i := range creditList {
		creditList[i].Credit /= multiple
	}
	return creditList, parseError
}

