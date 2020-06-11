package jwmodel

import (
	"bytes"
	"encoding/json"
	"github.com/Alex-zs/oucjwparser/util"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// 教务会话
type JwSession struct {
	UserCode   string `json:"userCode"`			// 学号
	PassWord   string `json:"password"`			// 密码
	JSESSIONID string `json:"j_session_id"`		// 教务会话ID
}

// 获取新的会话ID
func (session *JwSession) getNewJSessionID()  {
	util.SimpleDo(func(req *fasthttp.Request, resp *fasthttp.Response) {
		req.SetRequestURI(JwHost)
		req.Header.SetMethod(http.MethodGet)
		resp.SkipBody = true
		if err := fasthttp.DoTimeout(req, resp, util.DefaultTimeout); err != nil {
			util.Log("刷新会话ID失败", err.Error())
		}else {
			session.JSESSIONID = extractSession(resp.Header.PeekCookie("JSESSIONID"))
		}
	})
}

// 从cookie中提取session value
func extractSession(cookie []byte) string {
	pat := "JSESSIONID=[0-9A-Z]+"
	re, _ := regexp.Compile(pat)
	return string(re.Find(cookie)[11:])
}

// 验证会话有效性
func (session *JwSession) Validate() bool {
	valid := false
	util.SimpleDo(func(req *fasthttp.Request, resp *fasthttp.Response) {
		// 写入方法、cookie
		req.Header.SetMethod(fasthttp.MethodPost)
		req.SetRequestURI(OnlineMessage)
		req.Header.SetCookie("JSESSIONID", session.JSESSIONID)

		// 写入参数
		args := &fasthttp.Args{}
		args.Add("hidOption", "getOnlineMessage")
		args.WriteTo(req.BodyWriter())

		if err := fasthttp.DoTimeout(req, resp, util.DefaultTimeout); err != nil {
			util.Log("会话有效性验证失败", err.Error())
		}else {
			// 通过判断body是否为空验证有效性
			valid = len(resp.Body()) != 0
		}

	})
	return valid
}

var LoginSuccess = 0
var LoginFailure = 0

// 登录教务系统
func (session *JwSession) Login(userCode, passWord string) bool  {
	session.PassWord = passWord
	session.UserCode = userCode
	// 登录成功的标识
	loginSuccess := false

	// 获取新的会话ID
	session.getNewJSessionID()

	usercode := util.Base64Encoding([]byte(userCode + ";;" + session.JSESSIONID))
	result := 0
	for i := 0; i < len(passWord); i++ {
		charType := util.CharType(rune(passWord[i]))
		num := 0
		switch  charType{
		case util.Digit:
			num = 8
		case util.Lowercase:
			num = 4
		case util.Capital:
			num = 2
		case util.Others:
			num = 1
		}
		result = result | num
	}
	// 尝试识别验证码，最多尝试三次
	for tryTimes := 0; tryTimes < 3; tryTimes++ {
		// 获取验证码图片
		captchaPath := session.getCaptcha()
		// 识别验证码
		captcha := recognizeCaptcha(*captchaPath).Value
		password := util.MD5(util.MD5(passWord) + util.MD5(strings.ToLower(captcha)))
		data := ""
		util.SimpleDo(func(req *fasthttp.Request, resp *fasthttp.Response) {
			args := &fasthttp.Args{}
			args.Add("_u" + captcha, usercode)
			args.Add("_p" + captcha, password)
			args.Add("randnumber", captcha)
			args.Add("isPasswordPolicy", "1")
			args.Add("txt_mm_expression", strconv.Itoa(result))
			args.Add("txt_mm_length", string(len(passWord)))
			args.Add("txt_mm_userzh", "0")

			req.Header.SetRequestURI(Login)
			req.Header.SetMethod(fasthttp.MethodPost)
			req.Header.SetCookie("JSESSIONID", session.JSESSIONID)
			req.Header.SetReferer(JwHost)
			args.WriteTo(req.BodyWriter())

			err := fasthttp.DoTimeout(req, resp, util.DefaultTimeout)
			if err != nil {
				util.Log("请求登录失败", err.Error())
			}else {
				data = string(resp.Body())
				if strings.Index(data, "操作成功") > 0 {
					loginSuccess = true
					LoginSuccess += 1
				}else {
					LoginFailure += 1
				}
			}

			defer os.Remove(*captchaPath)
		})
		if loginSuccess {
			break
		}
	}
	return loginSuccess
}

// 尝试获取验证码图片
// 返回图片的本地路径,获取失败返回空字符串
func (session *JwSession) getCaptcha() *string {
	captchaPath := ""
	util.SimpleDo(func(req *fasthttp.Request, resp *fasthttp.Response) {
		req.SetRequestURI(Captcha)
		req.Header.SetCookie("JSESSIONID", session.JSESSIONID)

		err := fasthttp.DoTimeout(req, resp, util.DefaultTimeout)
		if err != nil {
			util.Log("获取验证码图片失败", err.Error())
			return
		}

		// 创建临时文件存放验证码图片
		captchaImg, err := ioutil.TempFile("", "*.jpeg")
		if err != nil {
			util.Log("创建图片验证码失败", err.Error())
			return
		}
		_, err = captchaImg.Write(resp.Body())
		if err != nil {
			util.Log("数据写入验证码图片失败", err.Error())
		}else {
			captchaPath = captchaImg.Name()
		}
	})
	return &captchaPath
}

// 识别验证码的接口
const RecognizeURL = "https://itstudio.club/ocr/jw"

// 识别接口返回的body
type CaptchaBody struct {
	SpendTime int `json:"spend_time(ms)"`	// 识别耗费的时间 毫秒
	TimeStamp string `json:"time"`			// 时间戳
	Value string `json:"value"`				// 识别值
}

// 开始识别验证码
func recognizeCaptcha(captchaPath string) *CaptchaBody {
	captchaBody := CaptchaBody{}

	// 新建缓冲区，用于存放图片
	bodyBuffer := &bytes.Buffer{}
	// 创建multipart文件写入器
	bodyWriter := multipart.NewWriter(bodyBuffer)
	// 写入表单
	fileWriter, err := bodyWriter.CreateFormFile("image", path.Base(captchaPath))
	if err != nil {
		util.Log(err.Error())
		return &captchaBody
	}
	file, _ := os.Open(captchaPath)
	defer file.Close()
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		util.Log(err.Error())
	}
	bodyWriter.Close()

	util.SimpleDo(func(req *fasthttp.Request, resp *fasthttp.Response) {
		req.Header.SetContentType(bodyWriter.FormDataContentType())
		req.SetBody(bodyBuffer.Bytes())
		req.Header.SetMethod(fasthttp.MethodPost)
		req.SetRequestURI(RecognizeURL)
		err := fasthttp.DoTimeout(req, resp, util.DefaultTimeout)

		if err != nil {
			util.Log("识别验证码失败", err.Error())
			return
		}

		json.Unmarshal(resp.Body(), &captchaBody)
	})

	return &captchaBody
}


